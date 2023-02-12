import { useState } from "react";

function SearchHistory({searchHistory}) {
    const [selectedSearchHistory, setSelectedSearchHistory] = useState(searchHistory)

    return (
        <div className="p-4 w-full">
            <div className="mb-16">
                <p className="w-full text-3xl border-b py-2 mb-4">Search History</p>
                <div className="cursor-pointer py-4 rounded-md mb-4 px-2 hover:bg-gray-100 dark:hover:bg-gray-400" onClick={()=>localStorage.removeItem('search')}>
                    <p className="text-lg">Clear Local Search History</p>
                    
                    <p className="text-sm text-gray-500">Remove Searches you have performed on this device</p>
                </div>
                {/* <div>
                    <label>
                        <input type='checkbox' className="h-4 w-4 mr-2" checked={selectedSearchHistory?'checked':''} onChange={()=>setSelectedSearchHistory(!selectedSearchHistory)} defaultValue=''></input>
                        Enable Search History
                    </label>
                </div> */}
            </div>
        </div>
    );
}

SearchHistory.defaultProps = {
    searchHistory: true
}

export default SearchHistory;