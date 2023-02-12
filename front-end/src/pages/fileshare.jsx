import { Fragment, useRef } from "react";
import AuthenticatedNavbar from "../components/authenticated-navbar";
import React, { useState, useEffect } from "react";
import axios from "axios";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faFile, faSearch, faShare, faXmark } from '@fortawesome/free-solid-svg-icons'
import Directory from "../components/directory";
import SearchDirectory from "../components/search-directory";
import useSWR from 'swr';
import PermissionForm from "../components/permissionForm";

function FileShare() {
	document.title = "File Share";

	if (localStorage.theme === 'dark') {
		document.documentElement.classList.add('dark')
	} else {
		document.documentElement.classList.remove('dark')
	} 

	// Handle Child Component References
	const directoryRef = useRef(null);

	const [isClicked, setClicked] = useState(false);

	// Folder directory data
	const [directory, setDirectory] = useState([]);
	const [currentFolder, setCurrentFolder] = useState("root");
	const [directoryMode, setDirectoryMode] = useState("root");

	const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]

	const fetcher = (url, token) =>
    axios
        .get(url, { headers: { "X-CSRF-TOKEN": csrfToken } })
        .then((res) => res.data);

    const { data: content, mutate: contentMutate } = useSWR(`/api/files/${currentFolder}`, fetcher);
    const { data: allContents, mutate: allContentMutate } = useSWR(`/api/contents`, fetcher);

	const [search, setSearch] = useState('')
	const [allSearch, setAllSearch] = useState([])
	const [formFocus, setFormFocus] = useState(false)
	const [searching, setSearching] = useState(false)



	// Create Folder Form methods
	const [createFolderForm, setCreateFolderForm] = useState(false);    // handles create folder modal state
	const [folderName, setFolderName] = useState("");					// stores new folder name
	const handleCreateFolderSwitch = () => {							// changes create new folder form visibility
		setCreateFolderForm(!createFolderForm);
		setFolderName("");												// delete folder name from form field
	};
	const handleCreateFolderForm = async() => {							// handles new folder creation
		try {
			axios.defaults.headers.post["X-CSRF-TOKEN"] = csrfToken
			const response = await axios({
				method: "post",
				url: "/api/create-folder",
				data: {
					Name: folderName,
					Parent: currentFolder
				},
				withCredentials: true
			});

			if (response.status === 201) {
				handleCreateFolderSwitch();									// close form
				contentMutate();						// update client data
			}
		} catch (error) {
			if (error.response.status === 500) {
				alert("something went wrong");
			} else {
				alert(error.response.data)
			};
		}
	};



	// Handle Content Clicks
	const [selectedContent, setSelectedContent] = useState(null);		// stores selected content data
	const handleContentClick = (data) => {
		setSelectedContent(data);
	};
	// Handle Directory traversal
	const switchDirectoryMode = mode => {							// handler to switch between directory modes
		setSelectedContent(null);									// clear selected content
		setDirectory([]);											// clear all directory data
		setCurrentFolder(mode);										// set current folder to directory mode value
		setDirectoryMode(mode);										// set directory mode
	}




	// Edit Object Form methods
	const [objectName, setObjectName] = useState("")
	const [editObjectNameForm, setEditObjectNameForm] = useState(false);    	// handles edit object modal state
	//const [folderName, setFolderName] = useState("");					// stores new folder name
	const handleEditObjectNameSwitch = () => {								// changes create new folder form visibility
		setEditObjectNameForm(!editObjectNameForm);
		setObjectName("")
	};
	const handleEditObjectNameForm = async() => {							// handles new folder creation
		try {
			const response = await axios({
				method: "post",
				url: `/api/edit/${selectedContent.type}/${selectedContent.details.ID}`,
				data: {
					Name: objectName
				},
				withCredentials: true
			});

			if (response.status === 201) {
				setObjectName("");											// delete folder name from form field
				setEditObjectNameForm(false)
				handleEditObjectNameSwitch();									// close form
				setSelectedContent(null)
				contentMutate();											// update client data
			}
		} catch (error) {
			if (error.response.status === 500) {
				alert("something went wrong");
			} else {
				alert(error.response.data)
			};
		}
	};



	// Edit Object Permissions Form methods
	const [editPermissionForm, setEditPermissionForm] = useState(false);	// handles edit object modal state
	const [target, setTarget] = useState("");								// stores target email
	const [targetPermission, setTargetPermission] = useState(0);			// stores target desired permission
	const handleEditPermissionSwitch = () => {								// changes edit permission form visibility
		setEditPermissionForm(!editPermissionForm);
		setTarget("");														// delete data in target field
		setTargetPermission(0);												// reset target permission field
	};
	const handleEditPermissionForm = async() => {							// handles new folder creation
		try {
			const response = await axios({
				method: "post",
				url: ("/api/share/"+selectedContent.type+"/"+selectedContent.details.ID),
				data: {
					Email: target,
					Permission: targetPermission
				},
				withCredentials: true
			});

			if (response.status === 201) {
				handleEditPermissionSwitch();									// close form
			}
		} catch (error) {
			if (error.response.status === 500) {
				alert("something went wrong");
			} else {
				alert(error.response.data)
			};
		}
	};



	// Reverse True & False Function
	function ReverseTrueFalse() {
		setClicked(!isClicked);
	};
	// Handle upload file state
	const [selectedFile, setSelectedFile] = useState();
	const [isFilePicked, setIsFilePicked] = useState(false);
	const changeHandler = (event) => {
		setSelectedFile(event.target.files[0]);
		setIsFilePicked(true);
	};

	// Customised changedHandler - SetClicked
	const changeHandler_show_button = (event) => {
		setClicked(true);
		ReverseTrueFalse();
	};

	// Handle File upload
	const handleSubmit = async (event) => {
		event.preventDefault();
		const formData = new FormData();
		// set file key to 'selectedFile'
		formData.append("selectedFile", selectedFile);
		try {
			axios.defaults.headers.post["X-CSRF-TOKEN"] = csrfToken
			const response = await axios({
				method: "post",
				url: ("/api/upload-file/"+currentFolder),
				data: formData,
				headers: { "Content-Type": "multipart/form-data" },
				withCredentials: true
			});

			if (response.status === 201) {
				contentMutate()
			}
		} catch (error) {
			if (error.response.status === 500) {
				alert("something went wrong");
			} else {
				alert(error.response.data)
			};
		}
	};

	useEffect (() => {
		if (localStorage.search !== undefined) {
		setAllSearch(JSON.parse(localStorage.search))
		}
	},[]);

	async function searchFiles() {
		var searchHistory = []
		if (localStorage.search !== undefined) {
			searchHistory = JSON.parse(localStorage.search)
		}
		if (!searchHistory.includes(search.trim()) && search.trim() !== '') {
		if (searchHistory.length >= 3) {
			searchHistory.pop()
		}
		searchHistory.unshift(search.trim())
		}
		localStorage.setItem('search', JSON.stringify(searchHistory))
		setAllSearch(JSON.parse(localStorage.search))
	}

	document.addEventListener('click', (event) => {
		setFormFocus(false)
	})

	return (
		<div className="flex flex-col justify-center items-center h-screen w-screen dark:bg-gray-800">
			<AuthenticatedNavbar />
			<div className="flex flex-col flex-grow container dark:bg-gray-800 dark:text-white mb-6">
				<div className="flex items-center justify-center mt-8">
					<form method="GET"  onSubmit={(e) => {e.preventDefault();searchFiles()}}>
						<div className="relative">
							<span className="absolute top-0 left-0 flex items-center pl-2">
								<button type="submit" className="p-1 focus:outline-none">
									<p className="text-black hover:text-gray-500">
										<FontAwesomeIcon icon={faSearch} className='' size="lg"/>
									</p>
								</button>
							</span>
							<div onClick={(e) => {e.stopPropagation(); setFormFocus(true)}}>
							<input type="search" name="search" className="border w-full border-black rounded-md pl-10 py-1 pr-2 w-[500px] text-black rounded-md text-gray-900 focus:outline-none" placeholder="Search..." autoComplete="off" value={search} onChange={e => {setSearch(e.target.value); setSearching(true); setSelectedContent(null)}} ></input>
								<div className={`list-none absolute ${formFocus ? '' : 'hidden'}`}>
									{allSearch.filter(item =>{
										return(item.includes(search.toLowerCase()))
									})
									.map((searchValue, index) => (
										<li key={index} className='cursor-pointer w-[500px] py-1 px-4 overflow-hidden text-ellipsis bg-white text-black hover:bg-gray-300' onClick={() => {setSearch(searchValue); setSearching(true); setSelectedContent(null)}} >{searchValue}</li>
									))}
									{/* {allContents.filter(item =>{
										const fullName = JSON.parse(item)['Document Name'].toLowerCase()
										return(fullName.includes(search.toLowerCase()))
									})
									.map((searchValue, index) => (
										<li key={index} className='cursor-pointer w-[500px] py-1 px-4 overflow-hidden text-ellipsis bg-white text-black hover:bg-gray-300' onClick={() =>{ setSearch(JSON.parse(searchValue)['Document Name']); searchHandler()}}>{JSON.parse(searchValue)['Document Name']}</li>
									))} */}
								</div>
							</div>
						</div>
					</form>
				</div>

				<div className="flex flex-row flex-grow w-full mt-6">
					{/* Directory Types */}
					<Fragment>
						<div className="flex flex-col justify-start items-start mr-6">
							<div 
								className="cursor-pointer overflow-hidden select-none flex flex-row justify-start items-center bg-transparent hover:bg-gray-300 hover:dark:bg-gray-600 text-gray-500 dark:text-gray-400 hover:text-gray-900 hover:dark:text-gray-300 font-semibold transition duration-150 ease-in-out rounded-xl p-3 mb-2 w-64"
								onClick={() => {switchDirectoryMode("root"); setSearching(false); setSearch("")}}
							>
								<div className="flex justify-center items-center w-[25px]">
									<FontAwesomeIcon icon={faFile} className='' size="lg"/>
								</div>
								<p className="pl-3 mr-1">My Files</p>
							</div>
							<div 
								className="cursor-pointer overflow-hidden select-none flex flex-row justify-start items-center bg-transparent hover:bg-gray-300 hover:dark:bg-gray-600 text-gray-500 dark:text-gray-400 hover:text-gray-900 hover:dark:text-gray-300 font-semibold transition duration-150 ease-in-out rounded-xl p-3 w-64"
								onClick={() => {switchDirectoryMode("shared"); setSearching(false); setSearch("")}}
							>
								<div className="flex justify-center items-center w-[25px]">
									<FontAwesomeIcon icon={faShare} className='' size="lg"/>
								</div>
								<p className="pl-3 mr-1">Shared With Me</p>
							</div>
						</div>

						<div className="bg-gray-400 w-0.5 my-2 dark:bg-gray-200" />
					</Fragment>

					<div className="flex flex-col w-full">
						{(content && content.Permission > 1 && !searching) &&
							<Fragment>
								<div className="ml-2">
									<button
										onClick={() => handleCreateFolderSwitch()}
										className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-3 px-20 rounded-full ml-8 mt-4 hover:shadow-2xl"
									>
										Create Folder
									</button>
									<button
										onClick={changeHandler_show_button}
										className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-3 px-20 rounded-full ml-8 my-4 hover:shadow-2xl"
									>
										Upload a File
									</button>
								</div>

								{/* File and Folder creation buttons */}
								{isClicked ? (
								<div>
									<p className="mt-4 ml-14">Please Select Your File Below :</p>
									{/* Grid 2 Columns */}
									<div className="grid grid-cols-2 gap-5">
										<form
											onSubmit={handleSubmit}
											name="selectedFile"
											encType="multipart/form-data"
											
										>
											<div className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-full ml-10 mb-2 w-80 hover:shadow-2xl">
												<input
													type="file"
													name="selectedFile"
													onChange={changeHandler}
													className="w-64 "
												/>
											</div>
											<div className="mt-2.5">
												<input
													type="submit"
													value="Upload"
													className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-full hover:shadow-2xl ml-10"
												/>
											</div>
										</form>
									</div>
								</div>
								) : (
								<p></p>
								)}
							</Fragment>
						}

						{ !searching ?
						<Directory 
							content={content}
							contentMutate={contentMutate}
							currentFolder={currentFolder}
							setCurrentFolder={setCurrentFolder}
							directory={directory}
							setDirectory={setDirectory}
							directoryMode={directoryMode}
							handleContentClick={handleContentClick} 
							selectedContent={selectedContent} 
							setSelectedContent={setSelectedContent} 
							handleEditPermissionSwitch={handleEditPermissionSwitch} 
							handleEditObjectNameSwitch={handleEditObjectNameSwitch}
						/>
						:
						<SearchDirectory 
							setSearching={setSearching}
							search={search}
							setSearch={setSearch}
							allContents={allContents}
							content={content}
							contentMutate={contentMutate}
							currentFolder={currentFolder}
							setCurrentFolder={setCurrentFolder}
							directory={directory}
							setDirectory={setDirectory}
							directoryMode={directoryMode}
							handleContentClick={handleContentClick} 
							selectedContent={selectedContent} 
							setSelectedContent={setSelectedContent} 
							handleEditPermissionSwitch={handleEditPermissionSwitch} 
							handleEditObjectNameSwitch={handleEditObjectNameSwitch}
						/>}
					</div>
				</div>
			</div>

			{/* Create form modal */}
			{createFolderForm &&
				<div className="fixed w-screen h-screen flex flex-col justify-center items-center z-10 bg-slate-900/50 inset-0">
					<div className="flex flex-col justify-center items-start bg-white text-gray-900 dark:bg-gray-800 dark:text-white rounded-2xl p-3">
						<div className="flex flex-row justify-between items-start w-full pb-3">
							<p className="text-lg font-semibold pl-1">New Folder</p>
							<FontAwesomeIcon 
								icon={faXmark} 
								size="lg" 
								className="cursor-pointer pt-0.5 pr-1.5" 
								onClick={() => handleCreateFolderSwitch()}
							/>
						</div>
						<input 
							type="text"
							id="FolderName"
							name="FolderName"
							className="text-base rounded-xl py-2 px-3 w-96 bg-gray-300 text-gray-500 dark:bg-gray-500 dark:text-gray-300 focus:bg-transparent transition duration-150 ease-in-out"
							placeholder="Folder Name"
							value={folderName}
							onChange={e => setFolderName(e.target.value)}
						/>
						<button 
							className="w-full bg-blue-500 rounded-xl py-2 mt-2"
							onClick={() => handleCreateFolderForm()}
						>
							<p className="cursor-pointer text-base text-white font-semibold">
								Create Folder
							</p>
						</button>
					</div>
				</div>
			}

			{/* Permission form modal */}
			{editPermissionForm &&
					<PermissionForm 
						handleEditPermissionSwitch={handleEditPermissionSwitch} 
						target={target}
						setTarget={setTarget}
						targetPermission={targetPermission}
						setTargetPermission={setTargetPermission}
						handleEditPermissionForm={handleEditPermissionForm}
						selectedContent={selectedContent}
					/>
			}

			{/* Object name form modal */}
			{editObjectNameForm &&
				<div className="fixed w-screen h-screen flex flex-col justify-center items-center z-10 bg-slate-900/50 inset-0">
					<div className="flex flex-col justify-center items-start bg-white text-gray-900 dark:bg-gray-800 dark:text-white rounded-2xl p-3">
						<div className="flex flex-row justify-between items-start w-full pb-3">
							<p className="text-lg font-semibold pl-1">Edit Object</p>
							<FontAwesomeIcon 
								icon={faXmark} 
								size="lg" 
								className="cursor-pointer pt-0.5 pr-1.5" 
								onClick={() => handleEditObjectNameSwitch()}
							/>
						</div>
						<div className="flex flex-row">
							<input 
								type="text"
								id="TargetEmail"
								name="TargetEmail"
								className="text-base rounded-xl py-2 px-3 w-96 bg-gray-300 text-gray-500 dark:bg-gray-500 dark:text-gray-300 focus:bg-transparent transition duration-150 ease-in-out"
								placeholder="Name"
								value={objectName}
								onChange={e => setObjectName(e.target.value)}
							/>
						</div>
						<button 
							className="w-full bg-blue-500 rounded-xl py-2 mt-2"
							onClick={() => handleEditObjectNameForm()}
						>
							<p className="cursor-pointer text-base text-white font-semibold">
								Edit
							</p>
						</button>
					</div>
				</div>
			}
		</div>
	);
}

export default FileShare;