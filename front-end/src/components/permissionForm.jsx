import React, { Fragment } from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faXmark } from '@fortawesome/free-solid-svg-icons';
import axios from 'axios';
import useSWR from 'swr';

const PermissionForm = ({ handleEditPermissionSwitch, target, setTarget, targetPermission, setTargetPermission, handleEditPermissionForm, selectedContent }) => {

    const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]

	const fetcher = (url, token) =>
    axios
        .get(url, { headers: { "X-CSRF-TOKEN": csrfToken } })
        .then((res) => res.data);

    const { data: content } = useSWR(`/api/${selectedContent.type}/permissions/${selectedContent.details.ID}`, fetcher);

    

    return (
        <div className="fixed w-screen h-screen flex flex-col justify-center items-center z-10 bg-slate-900/50 inset-0">
            <div className="flex flex-col justify-center items-start bg-white text-gray-900 dark:bg-gray-800 dark:text-white rounded-2xl p-3">
                <div className="flex flex-row justify-between items-start w-full pb-3">
                    <p className="text-lg font-semibold pl-1">Object Permissions</p>
                    <FontAwesomeIcon 
                        icon={faXmark} 
                        size="lg" 
                        className="cursor-pointer pt-0.5 pr-1.5" 
                        onClick={() => handleEditPermissionSwitch()}
                    />
                </div>
                {content != null &&
                    <Fragment>
                        <table className='table-auto w-full border border-gray-300 rounded-lg mb-3'>
                            <thead>
                                <tr className='text-gray-900 dark:text-white'>
                                    <th className='p-2'>User Email</th>
                                    <th className='p-2'>Permissions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {content.map((individual) => (
                                    <tr key={individual.ID} className='bg-transparent hover:bg-gray-200 text-gray-500 hover:text-gray-900 dark:text-gray-200 dark:hover:text-white dark:hover:bg-gray-700'>
                                        <td className='text-start p-2'>{individual.Email}</td>
                                        <td className='text-center p-2'>{individual.Permission === 1 && "Viewing"}{individual.Permission === 2 && "Editing"}{individual.Permission === 3 && "Administrator"}{individual.Permission === 4 && "Owner"}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </Fragment>
                }
                <div className="flex flex-row">
                    <input 
                        type="text"
                        id="TargetEmail"
                        name="TargetEmail"
                        className="text-base rounded-xl py-2 px-3 w-96 bg-gray-300 text-gray-500 dark:bg-gray-500 dark:text-white focus:bg-transparent dark:focus:bg-transparent transition duration-150 ease-in-out"
                        placeholder="Add Emails"
                        value={target}
                        onChange={e => setTarget(e.target.value)}
                    />
                    <select
                        className="text-base rounded-xl py-2 px-3 ml-2 bg-gray-300 text-gray-500 dark:bg-gray-500 dark:text-gray-300 focus:bg-transparent border-0 focus:border-2 border-gray-900 transition duration-150 ease-in-out"
                        value={targetPermission}
                        onChange={e => setTargetPermission(parseInt(e.target.value))}
                    >
                        <option value={3}>Administrator</option>
                        <option value={2}>Editor</option>
                        <option value={1}>Viewer</option>
                        <option value={0}>No Permission</option>
                    </select>
                </div>
                <button 
                    className="w-full bg-blue-500 rounded-xl py-2 mt-2"
                    onClick={() => handleEditPermissionForm()}
                >
                    <p className="cursor-pointer text-base text-white font-semibold">
                        Share
                    </p>
                </button>
            </div>
        </div>
    );
};

export default PermissionForm;