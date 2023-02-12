function Feedback() {
    return (
        <div className="p-4 w-full">
            <div className="mb-8">
                <p className="w-full text-3xl border-b py-2 mb-4">Feedback</p>
                <div>
                    <form>
                        <input type='text' placeholder='Title' className='border border-2 border-black w-full py-2 px-3 mb-4 text-gray-700 mb-3 focus:outline-none focus:border-2'></input>
                        <textarea className='border border-2 border-black w-full h-64 py-2 px-3 mb-4 text-gray-700 mb-3 focus:outline-none focus:border-2'></textarea>
                        <button type='submit' className='bg-black w-full py-2 px-3 text-white mb-3 focus:outline-none focus:border-2 hover:text-gray-200 dark:bg-gray-300 dark:text-black dark:hover:bg-gray-400'>Send</button>
                    </form>
                </div>
            </div>
        </div>
    );
}

export default Feedback;