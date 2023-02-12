import { Fragment, useState, useEffect } from "react";
import axios from 'axios'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faPaintBrush, faBell, faUser, faLock, faCreditCard, faSearch, faMessage, faFileContract, faShield } from '@fortawesome/free-solid-svg-icons'
import Navbar from "../components/navbar";
import AuthenticatedNavbar from "../components/authenticated-navbar";
import Account from "../components/settings/Account";
import Appearance from "../components/settings/Appearance";
import Notifications from "../components/settings/Notifications";
import PasswordAuthentication from "../components/settings/PasswordAuthentication";
import BillingPlans from "../components/settings/BillingPlans";
import SearchHistory from "../components/settings/SearchHistory";
import Feedback from "../components/settings/Feedback";

function Settings() {
    const [page, setPage] = useState('Account')
    const [theme, setTheme] = useState(localStorage.theme)
    const [username, setUsername] = useState('')

    document.title = 'Settings'

    useEffect (() => {
        const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]
        axios.defaults.headers.get["X-CSRF-TOKEN"] = csrfToken
        axios.get('/api/username')
        .then((res) => setUsername(res.data))
        .catch(err=> setUsername('')) 

        if (username === '') {
            setPage('Appearance')
        }
    
        if (localStorage.theme === 'dark') {
            document.documentElement.classList.add('dark')
        } else {
            document.documentElement.classList.remove('dark')
        } 
    },[theme]);

    function toggleTheme (selectedTheme) { 
        setTheme(selectedTheme)
        localStorage.setItem('theme', selectedTheme)
    }

    return (
        <Fragment>
            {username === '' ?
                <Navbar/>
            :
                <AuthenticatedNavbar username={username}/>
            }
            <div className="flex min-h-screen">
                <div className="border-r-2 border-black w-[335px] pl-4 py-4 dark:bg-gray-500">
                    <ul>
                        {username !== '' ?
                            <li className={`cursor-pointer py-1 px-2 rounded-l-md hover:bg-gray-200 flex ${username != ''?"block":"hidden"}`} onClick={() => {setPage('Account')}}>
                                <div className="w-5 mr-2">
                                    <FontAwesomeIcon icon={faUser} />
                                </div>
                                Account
                            </li>
                        :
                            <></>
                        }
                        <li className="cursor-pointer py-1 px-2 rounded-l-md hover:bg-gray-200 flex" onClick={() => {setPage('Appearance')}}>
                            <div className="w-5 mr-2">
                                <FontAwesomeIcon icon={faPaintBrush} />
                            </div>
                            Appearance
                        </li>
                        {username !== '' ?
                            <Fragment>
                                <li className="cursor-pointer py-1 px-2 rounded-l-md hover:bg-gray-200 flex hidden" onClick={() => {setPage('Notifications')}}>
                                    <div className="w-5 mr-2">
                                        <FontAwesomeIcon icon={faBell} />
                                    </div>
                                    Notifications
                                </li>
                                <li className="cursor-pointer py-1 px-2 rounded-l-md hover:bg-gray-200 flex" onClick={() => {setPage('PasswordAuthentication')}}>
                                    <div className="w-5 mr-2">
                                        <FontAwesomeIcon icon={faLock} />
                                    </div>
                                    Password and Authentication
                                </li>
                                {/* <li className="cursor-pointer py-1 px-2 rounded-l-md hover:bg-gray-200 flex" onClick={() => {setPage('BillingPlans')}}>
                                    <div className="w-5 mr-2">
                                        <FontAwesomeIcon icon={faCreditCard} />
                                    </div>
                                    Billing and Plans
                                </li> */}
                                <li className="cursor-pointer py-1 px-2 rounded-l-md hover:bg-gray-200 flex" onClick={() => {setPage('SearchHistory')}}>
                                    <div className="w-5 mr-2">
                                        <FontAwesomeIcon icon={faSearch} />
                                    </div>
                                    Search History
                                </li>
                            </Fragment>
                        :
                            <></>
                        }
                        {/* <li className="cursor-pointer py-1 px-2 rounded-l-md hover:bg-gray-200 flex" onClick={() => {setPage('Feedback')}}>
                            <div className="w-5 mr-2">
                                <FontAwesomeIcon icon={faMessage} />
                            </div>
                            Feedback
                        </li> */}
                        {/* <li className="cursor-pointer py-1 px-2 rounded-l-md hover:bg-gray-200 flex" onClick={() => {setPage('T&C')}}>
                            <div className="w-5 mr-2">
                                <FontAwesomeIcon icon={faFileContract} />
                            </div>
                            Terms and Conditions
                        </li> */}
                        {/* <li className="cursor-pointer py-1 px-2 rounded-l-md hover:bg-gray-200 flex" onClick={() => {setPage('PrivacyPolicy')}}>
                            <div className="w-5 mr-2">
                                <FontAwesomeIcon icon={faShield} />
                            </div>
                            Privacy Policy
                        </li> */}
                    </ul>
                </div>
                <div className="w-full dark:bg-gray-800 dark:text-white  dark:border-white">
                    { page === 'Account' ?
                        <Account />
                    :
                        <></>
                    }
                    { page === 'Appearance' ?
                        <Appearance  toggleTheme={toggleTheme} theme={localStorage.theme}/>
                    :
                        <></>
                    }
                    { page === 'Notifications' ?
                        <Notifications />
                    :
                        <></>
                    }
                    {page === 'PasswordAuthentication' ?
                        <PasswordAuthentication />
                    :
                        <></>
                    }
                    { page === 'BillingPlans' ?
                        <BillingPlans />
                    :
                        <></>
                    }
                    { page === 'SearchHistory' ?
                        <SearchHistory />
                    :
                        <></>
                    }
                    { page === 'Feedback' ?
                        <Feedback />
                    :
                        <></>
                    }
                    { page === 'T&C' ?
                        <Feedback />
                    :
                        <></>
                    }
                    { page === 'PrivacyPolicy' ?
                        <Feedback />
                    :
                        <></>
                    }
                </div>
            </div>
        </Fragment>
    );
}

export default Settings;