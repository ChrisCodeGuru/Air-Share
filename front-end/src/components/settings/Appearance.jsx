import { useState } from "react";

function Appearance({theme, maxFiles, toggleTheme}) {
    const [selectedMaxFiles, setSelectedMaxFiles] = useState(maxFiles)

    return (
        <div className="p-4 w-full">
            <div className="mb-8">
                <p className="w-full text-3xl border-b py-2 mb-4">Theme Preferences</p>
                <form className="flex p-2">
                    <div className="cursor-pointer w-[400px] rounded-md py-2 px-3 hover:bg-gray-100 dark:hover:bg-gray-400" onClick={e => toggleTheme('light')}>
                        <input type='radio' checked={theme!='dark'?'checked':''} onChange={e => toggleTheme('light')}>
                        </input>
                        <span>&nbsp;Light</span>
                        <div className="border-2 border-black rounded-md bg-white my-4 p-2 h-32 text-center">
                            <div className="flex pb-2 border-b justify-between">
                                <div className="bg-black h-3 rounded-md w-16"></div>
                                <div className="flex">
                                    <div className="bg-black h-3 rounded-md w-8 mr-3"></div>
                                    <div className="bg-black h-3 rounded-md w-8 mr-3"></div>
                                    <div className="bg-black h-3 rounded-md w-8"></div>
                                </div>
                            </div>
                            <p className="pt-4 font-bold dark:text-black">AIRSHARE</p>
                        </div>
                    </div>
                    <div className="cursor-pointer w-[400px] rounded-md py-2 px-3 hover:bg-gray-100 dark:hover:bg-gray-400" onClick={e => toggleTheme('dark')}>
                        <input type='radio' checked={theme==='dark'?'checked':''} onChange={e => toggleTheme('dark')}></input>
                        <span>&nbsp;Dark</span>
                        <div className="border-2 border-white rounded-md bg-black my-4 p-2 h-32 text-center">
                            <div className="flex pb-2 border-b justify-between">
                                <div className="bg-white h-3 rounded-md w-16"></div>
                                <div className="flex">
                                    <div className="bg-white h-3 rounded-md w-8 mr-3"></div>
                                    <div className="bg-white h-3 rounded-md w-8 mr-3"></div>
                                    <div className="bg-white h-3 rounded-md w-8"></div>
                                </div>
                            </div>
                            <p className="py-4 font-bold text-white">AIRSHARE</p>
                            <div className="flex justify-center">
                                <div className="bg-gray-700 rounded-t-md h-9 w-[75%]"></div>
                            </div>
                        </div>
                    </div>
                </form>
            </div>
            
            {/* <div className="mb-4">
                <p className="w-full text-3xl border-b py-2 mb-4">Maximum Page Size</p>
                <div>
                    Show &nbsp;
                    <select className='border-[1px] border-black rounded-md px-2 py-0.5 w-16 text-black' value={selectedMaxFiles} onChange={(e)=>setSelectedMaxFiles(e.target.value)} >
                        <option value="10">10</option>
                        <option value="15">15</option>
                        <option value="20">20</option>
                        <option value="25">25</option>
                    </select>
                    &nbsp; files per page
                </div>
            </div> */}
        </div>
    );
}

Appearance.defaultProps = {
    theme: 'light',
    maxFiles: 10
}

export default Appearance;