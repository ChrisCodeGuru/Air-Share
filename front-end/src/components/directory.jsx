import React, { Fragment } from 'react';
import axios from 'axios';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faChevronLeft, faArrowRight, faFolder, faFile, faArrowUpFromBracket } from '@fortawesome/free-solid-svg-icons';

const Directory = ({ 
        content,
        contentMutate,
        setCurrentFolder, 
        directory, 
        setDirectory, 
        directoryMode, 
        handleContentClick, 
        selectedContent, 
        setSelectedContent, 
        handleEditPermissionSwitch, 
        handleEditObjectNameSwitch
    }) => {

    const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]

    // Handle folder double click
    const handleFolderClick = (event, data) => {
		if (event.detail === 2) {
			setSelectedContent(null);									// remove selected content
			setCurrentFolder(data.details.ID);							// changes current folder
			setDirectory(directory => [...directory, data.details]);	// appends to directory
		} else {
			handleContentClick(data);
		};
	};

    // Handle Folder Delete
	const handleFolderDelete = async (ID) => {
        try {
            axios.defaults.headers.post["X-CSRF-TOKEN"] = csrfToken
            const response = await axios({
                method: "POST",
                url: ("/api/delete-folder/" + ID),
                withCredentials: true
            });

            if (response.status === 200) {
                alert("successfully deleted");
                contentMutate();
                setSelectedContent(null);
            }
        } catch (error) {
            if (error.response.status === 500) {
                alert("something went wrong");
            } else {
                alert(error.response.data)
            };
        }
	};

    // Handle File download
    const handleFileDownload = async (ID) => {
        try {
            axios.defaults.headers.get["X-CSRF-TOKEN"] = csrfToken
            const response = await axios({
                method: "GET",
                url: ("/api/download-file/" + ID),
                withCredentials: true,
                responseType: 'blob'
            });

            if (response.status === 200) {
                const url = window.URL.createObjectURL(new Blob([response.data]));
                const link = document.createElement('a');
                link.href = url;
                link.setAttribute('download', response.headers['content-disposition'].split('filename=')[1].split('.')[0] + "." + response.headers['content-disposition'].split('.')[1].split(';')[0]);
                document.body.appendChild(link);
                link.click();
            }
        } catch (error) {
            if (error.response.status === 500) {
                alert("something went wrong");
            } else {
                alert(error.response.data)
            };
        }
	}

    // Handle File delete
    const handleFileDelete = async (ID) => {
        try {
            axios.defaults.headers.post["X-CSRF-TOKEN"] = csrfToken
            const response = await axios({
                method: "POST",
                url: ("/api/delete-file/" + ID),
                withCredentials: true
            });

            if (response.status === 200) {
                alert("successfully deleted");
                contentMutate();
                setSelectedContent(null);
            }
        } catch (error) {
			if (error.response.status === 500) {
				alert("something went wrong");
			} else {
				alert(error.response.data)
			};
		}
	};

    // Handle Directory Traversal
	const handleDirectoryBack = () => {									// handles traverse previous directory
		if (directory.length === 1) {
			setCurrentFolder(directoryMode);							// set to base directory taken from directory mode if only one folder
		} else {
			setCurrentFolder(directory[directory.length-2].ID);			// set UUID of 2nd order folder parent if length more than 1
		};

		setSelectedContent(null);										// remove selected content
		directory.splice(-1);											// remove last element of directory
		setDirectory(directory);										// update directory state
	};

    return (
        <div className='flex flex-row grow w-full'>
            <div className="flex flex-col flex-grow justify-start items-start mx-6">
                {/* Directory */}
                {directory.length !== 0 &&
                    <div className="flex flex-row justify-start items-center">
                        <div 
                            className="cursor-pointer flex justify-center items-center w-12 h-12 rounded-xl bg-transparent hover:bg-gray-300 hover:dark:bg-gray-600 text-gray-500 hover:text-gray-900 hover:dark:text-gray-300 font-semibold transition duration-150 ease-in-out"
                            onClick={() => handleDirectoryBack()}
                        >
                            <FontAwesomeIcon icon={faChevronLeft} className='' size="lg"/>
                        </div>
                        {directory.map((directory, index) => 
                            <Fragment>
                                {index !== 0 &&
                                    <FontAwesomeIcon icon={faArrowRight} className='text-gray-500' size="lg"/>
                                }
                                <div 
                                    className="cursor-pointer select-none flex flex-row justify-center items-center p-3 mx-2 rounded-xl bg-transparent hover:bg-gray-300 hover:dark:bg-gray-600 text-gray-500 hover:text-gray-900 hover:dark:text-gray-300 font-semibold transition duration-150 ease-in-out"
                                    onClick={() => {}}
                                >
                                    <FontAwesomeIcon icon={faFolder} className='' size="lg"/>
                                    <p className="pl-3">{directory.Name}</p>
                                </div>
                            </Fragment>
                        )}
                    </div>
                }

                {/* Directory contents */}
                <div className="flex flex-row">
                    {content &&
                        <Fragment>
                            {content.Folders !== null &&
                                content.Folders.map((folder, index) => 
                                    <div 
                                        className="cursor-pointer flex flex-row justify-center items-center select-none p-3 m-2 border rounded-xl" 
                                        onClick={e => handleFolderClick(e, {type: "folder", details: folder})}
                                    >
                                        <FontAwesomeIcon icon={faFolder} className='' size="lg"/>
                                        <p className="pl-3">{folder.Name}</p>
                                    </div>
                                )
                            }
                            {content.Files !== null &&
                                content.Files.map((file, index) => 
                                    <div 
                                        className="cursor-pointer flex flex-row justify-center items-center select-none p-3 m-2 border rounded-xl" 
                                        onClick={() => handleContentClick({type: "file", details: file})}
                                    >
                                        <FontAwesomeIcon icon={faFile} className='' size="lg"/>
                                        <p className="pl-3">{file.Name}</p>
                                    </div>
                                )
                            }
                        </Fragment>
                    }
                </div>
            </div>

            {/* Selected File Details */}
            {selectedContent !== null &&
                <div className='flex flex-row h-full pt-4'>
                    <div className="bg-gray-400 w-0.5 my-2" />
                    <div className="flex flex-col justify-start items-center ml-6 w-64">
                        <p className="text-2xl font-bold text-gray-900 pb-8 pt-4">Selected File</p>
                        {selectedContent.type === "folder" &&
                            <FontAwesomeIcon icon={faFolder} className='text-gray-900dark:text-white' size="10x"/>
                        }
                        {selectedContent.type === "file" &&
                            <FontAwesomeIcon icon={faFile} className='' size="10x"/>
                        }
                        <p className="text-gray-900 text-2xl mb-4">{selectedContent.details.Name}</p>
                        {selectedContent.type === "folder" && 
                            <p className="text-gray-900 text-base text-start mb-8 w-full"><span className="font-bold">Permission level:</span> {selectedContent.details.Permission === 1 && "Viewing"}{selectedContent.details.Permission === 2 && "Editing"}{selectedContent.details.Permission === 3 && "Administrator"}{selectedContent.details.Permission === 4 && "Owner"}</p>
                        }
                        {selectedContent.type === "file" &&
                            <Fragment>
                                <p className="text-gray-900 text-base text-start mb-2 w-full"><span className="font-bold">Permission level:</span> {selectedContent.details.Permission === 1 && "Viewing"}{selectedContent.details.Permission === 2 && "Editing"}{selectedContent.details.Permission === 3 && "Administrator"}{selectedContent.details.Permission === 4 && "Owner"}</p>
                                <p className="text-gray-900 text-base text-start mb-2 w-full"><span className="font-bold">Sensitive:</span> {selectedContent.details.Sensitive === "t" && "True"}{selectedContent.details.Sensitive === "f" && "False"}</p>
                                <p className="text-gray-900 text-base text-start mb-8 w-full break-words"><span className="font-bold">SHA256 Checksum:</span> {selectedContent.details.Hash}</p>
                            </Fragment>
                        }
                        {selectedContent.details.Permission >= 2 &&
                            <div className="flex flex-row justify-between items-center w-full mb-2">
                                {selectedContent.type === "folder" &&
                                    <div 
                                        className="cursor-pointer flex justify-center items-center flex-grow bg-blue-500 p-3 rounded-xl"
                                        onClick={() => handleEditObjectNameSwitch()}
                                    >
                                        <p className="font-semibold text-white">Edit</p>
                                    </div>
                                }
                                {selectedContent.type === "file" &&
                                    <div 
                                    className="cursor-pointer flex justify-center items-center flex-grow bg-blue-500 p-3 rounded-xl"
                                    onClick={() => handleFileDownload(selectedContent.details.ID)}
                                    >
                                        <p className="font-semibold text-white">Download</p>
                                    </div>
                                }
                                {selectedContent.details.Permission >= 3 &&
                                    <div 
                                        className="cursor-pointer flex justify-center items-center w-12 h-12 ml-2 text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-200 rounded-xl border border-gray-200 transition duration-150 ease-in-out"
                                        onClick={() => handleEditPermissionSwitch()}
                                    >
                                        <FontAwesomeIcon icon={faArrowUpFromBracket} className='' size=""/>
                                    </div>
                                }
                            </div>
                        }
                        {selectedContent.details.Permission === 4 &&
                            <div 
                                className="cursor-pointer flex justify-center items-center bg-red-500 w-full p-3 rounded-xl"
                                onClick={() => {
                                    if (selectedContent.type === "folder") {
                                        handleFolderDelete(selectedContent.details.ID)
                                    } else if (selectedContent.type === "file") {
                                        handleFileDelete(selectedContent.details.ID)
                                    }
                                }}
                            >
                                <p className="font-semibold text-white">Delete</p>
                            </div>
                        }
                    </div>
                </div>
            }
        </div>
    )
};

export default Directory;