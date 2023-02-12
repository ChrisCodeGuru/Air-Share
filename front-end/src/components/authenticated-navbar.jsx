import { Link } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faGear, faBell } from '@fortawesome/free-solid-svg-icons'
import { useState, useEffect } from "react";
import axios from "axios";


function AuthenticatedNavbar() {
    const [dropdown, setDropdown] = useState(false)
    const [username, setUsername] = useState("")
    const SetDropdown = () => {
        if (dropdown === false){
            setDropdown(true)
        }
        else {
            setDropdown(false)
        }
    }

    const handleLogout = async() => {
        const response = await axios({
			method: "get",
			url: "/api/logout",
			withCredentials: true
		});

        if (response.status === 200) {
            localStorage.removeItem("search")
            window.location.replace("/");
        }
    }

    useEffect (() => {
        const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]
        axios.defaults.headers.get["X-CSRF-TOKEN"] = csrfToken
        axios.get('/api/username')
        .then((res) => {
            if (res.status === 200){
                setUsername(res.data)
            }
        })
        .catch(err=> setUsername(''))
    },[]);

    return (
        <div className="py-4 px-8 w-full flex items-center justify-between border-b dark:bg-gray-900">
            <div>
                <Link to='/' className="text-2xl font-bold tracking-wider mr-8 dark:text-white">AIRSHARE</Link>
            </div>
            <div className="flex justify-between w-[600px]">
                <div className="px-2 w-[250px] flex items-center justify-between text-slate-400 font-bold">
                    <Link to='/fileshare' className="hover:text-slate-500 dark:hover:text-slate-300">My Files</Link>
                    {/* <Link to='#' className="hover:text-slate-500 dark:hover:text-slate-300">xxxxx</Link>
                    <Link to='#' className="hover:text-slate-500 dark:hover:text-slate-300">xxxxx</Link> */}
                </div>

                <div className="px-2 ml-8 max-w-[300px] font-bold">
                    <ul className="flex">
                        <li>
                            <p className="w-[150px] overflow-hidden text-ellipsis mr-12 text-center">
                                <Link to='#' className="text-black p-2 hover:underline dark:text-white dark:decoration-white" onClick={() => {SetDropdown()}}>{username}</Link>
                            </p>
                            {dropdown ?
                                <div className='block p-2 pl-2 mt-1 bg-white rounded-md border-2 absolute right-[170px]'>
                                    <ul className="space-y-2 text-gray-600">
                                        {/* <li>
                                            <Link href="#" className="flex p-2 font-medium rounded-md hover:text-black hover:bg-gray-100">Manage Account</Link>
                                        </li> */}
                                        <li>
                                            <div onClick={() => handleLogout()} className="cursor-pointer flex p-2 font-medium rounded-md hover:text-black hover:bg-gray-100">Sign Out</div>
                                        </li>
                                    </ul>
                                </div>
                            :
                                <></>
                            }
                        </li>
                        <li>
                            {/* <Link to='#'>
                                <div className='flex px-2 mr-8 hover:text-gray-600 cursor-pointer'>
                                    <p className="dark:text-white dark:hover:text-gray-300">
                                        <FontAwesomeIcon icon={faBell} size="lg"/>
                                    </p>
                                    { notifications > 0 ?
                                        <span className="relative bottom-0.5 right-1.5 bg-white p-1 rounded-full h-3 w-3 flex justify-center items-center dark:bg-black">
                                            <span className=" flex rounded-full h-2 w-2 bg-sky-500 absolute justify-center items-center">
                                                <span className="rounded-full h-2 w-2 bg-sky-500 animate-ping"></span>       
                                            </span>
                                        </span>
                                    :
                                        <span className="relative bottom-0.5 right-1.5 p-1 h-3 w-3"></span>
                                    }
                                </div>
                            </Link> */}
                            {/* <Link href='#'>
                                <div className='cursor-pointer flex px-2 hover:text-gray-300'>
                                    <p className='cursor-pointer'><FontAwesomeIcon icon={faBell} size="lg" shake/></p>
                                    <span className='relative bottom-1.5 text-xs'>{notifications}</span>
                                </div>
                            </Link> */}
                        </li>
                        <li>
                            <div className="hover:text-gray-600">
                                <Link to='/settings'>
                                    <p className="dark:text-white dark:hover:text-gray-300">
                                        <FontAwesomeIcon icon={faGear} className='fa-hover-gray-400' size="lg"/>
                                    </p>
                                </Link>
                            </div>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    );
}

// AuthenticatedNavbar.defaultProps = {
//     notifications: 9
// }

export default AuthenticatedNavbar;