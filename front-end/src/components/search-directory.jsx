import React, { Fragment } from 'react';
import axios from 'axios';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faChevronLeft, faArrowRight, faFolder, faFile, faArrowUpFromBracket } from '@fortawesome/free-solid-svg-icons';

const SearchDirectory = ({ 
        setSearching,
        search,
        setSearch,
        allContents,
        content,
        contentMutate,
        setCurrentFolder, 
        setDirectory, 
        handleContentClick, 
        selectedContent, 
        setSelectedContent, 
        handleEditPermissionSwitch, 
        handleEditObjectNameSwitch
    }) => {

    const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]

    // Handle File Traverse
    const handleFolderClick = (data) => {
        var dirList = []
        if (data.details=="") {
            if (data.permission = 4){
                setCurrentFolder("root");							                    // changes to root folder
            } else {
			    setCurrentFolder("root");							                    // changes to shared folder   
            }
            setSelectedContent(null);									                // remove selected content
			setDirectory([]);	                                                        // remove all directory
        }
		else {
            for (let i = 0; i < data.details.length; i++) {
                dirList.push(JSON.parse(data.details[i]))
            }

			setSelectedContent(null);									            // remove selected content
			setCurrentFolder(dirList[data.details.length-1].ID);					// changes current folder
			setDirectory(dirList);	                                                // appends to directory
		}
        setSearch("")
        setSearching(false)
    }

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

    return (
        <div className='flex flex-row grow w-full'>
            <div className="flex flex-col flex-grow justify-start items-start mx-6">
                {/* Directory contents */}
                <div className="flex flex-row">
                    {content &&
                        <Fragment>
                            {allContents !== null &&
                                allContents
                                .filter(item =>{
                                    const fullName = JSON.parse(item)['Document Name'].toLowerCase()
                                    return(fullName.includes(search.toLowerCase()))
                                })
                                .map((file, index) => 
                                    <div 
                                        className="cursor-pointer flex flex-row justify-center items-center select-none p-3 m-2 border rounded-xl" 
                                        onClick={e => handleContentClick({type: "file", details: {"ID": JSON.parse(file)['File ID'], Name: JSON.parse(file)['Document Name'], Owner: JSON.parse(file)['User Email'], Permission: JSON.parse(file)['Permission'], Path: JSON.parse(file)['File Path']}, Folder: JSON.parse(file)['File Directory']})}
                                    >
                                        {JSON.parse(file)['Type'] == "File" ?
                                            <FontAwesomeIcon icon={faFile} className='' size="lg"/>
                                        :
                                            <FontAwesomeIcon icon={faFolder} className='' size="lg"/>
                                        }
                                        <p className="pl-3">{JSON.parse(file)['Document Name']}</p>
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
                            <FontAwesomeIcon icon={faFolder} className='text-gray-900' size="10x"/>
                        }
                        {selectedContent.type === "file" &&
                            <FontAwesomeIcon icon={faFile} className='' size="10x"/>
                        }
                        <p className="text-gray-900 dark:text-white text-2xl mb-4">{selectedContent.details.Name}</p>
                        <p className="text-gray-900 dark:text-white text-base text-start mb-2 w-full"><span className="font-bold">Owner:</span> {selectedContent.details.Owner}</p>
                        <p className="text-gray-900 dark:text-white text-base text-start mb-2 w-full" ><span className="font-bold">Permission level: </span> {selectedContent.details.Permission === 1 && "Viewing"}{selectedContent.details.Permission === 2 && "Editing"}{selectedContent.details.Permission === 3 && "Administrator"}{selectedContent.details.Permission === 4 && "Owner"}</p>
                        <p className="text-gray-900 dark:text-white text-base text-start mb-8 w-full cursor-pointer" onClick={e => handleFolderClick({type: "folder", permission: selectedContent.details.Permission, details: selectedContent.Folder})}>
                            <span className="font-bold">File Path:</span>
                            {selectedContent.details.Permission == 4 ?
                                (selectedContent.details.Path).startsWith("/")?
                                "My Files" + selectedContent.details.Path
                                :
                                "My Files/" + selectedContent.details.Path
                            :
                                (selectedContent.details.Path).startsWith("/")?
                                "Shared With Me" + selectedContent.details.Path
                                :
                                "Shared With Me/" + selectedContent.details.Path
                            }
                        </p>
                        {selectedContent.details.Permission >= 2 &&
                            <div className="flex flex-row justify-between items-center w-full mb-2">
                                <div 
                                    className="cursor-pointer flex justify-center items-center flex-grow bg-blue-500 p-3 rounded-xl"
                                    onClick={() => handleEditObjectNameSwitch()}
                                >
                                    <p className="font-semibold text-white">Edit</p>
                                </div>
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
                        {selectedContent.type === "file" &&
                            <div 
                            className="cursor-pointer flex justify-center items-center bg-blue-500 p-3 rounded-xl mb-2 w-full"
                            onClick={() => handleFileDownload(selectedContent.details.ID)}
                            >
                                <p className="font-semibold text-white">Download</p>
                            </div>
                        }
                        {selectedContent.details.Permission === 4 &&
                            <div 
                                className="cursor-pointer flex justify-center items-center bg-red-500 w-full p-3 rounded-xl"
                                onClick={() => {
                                    if (selectedContent.type === "folder") {
                                        //handleFolderDelete(selectedContent.details.ID)
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

export default SearchDirectory;