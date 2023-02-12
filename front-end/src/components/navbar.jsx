import { Link } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faGear } from '@fortawesome/free-solid-svg-icons'

function Navbar() {
  return (
    <div className="py-4 px-8 w-full flex items-center justify-between border-b dark:bg-gray-900">
        <div>
            <Link to='/' className="text-2xl font-bold tracking-wider mr-8 dark:text-white">AIRSHARE</Link>
        </div>
        {/* <div className="flex justify-between w-[600px]"> */}
            {/* <div className="px-2 w-[250px] flex items-center justify-between text-slate-400 font-bold">
                <Link to='#' className="hover:text-slate-500 dark:hover:text-slate-300">Solutions</Link>
                <Link to='#' className="hover:text-slate-500 dark:hover:text-slate-300">About</Link>
                <Link to='#' className="hover:text-slate-500 dark:hover:text-slate-300">Pricing</Link>
            </div> */}
            <div className="px-2 ml-8 w-[300px] flex items-center justify-between font-bold">
                <Link to='/login' className="px-2 py-3 rounded-md text-black hover:bg-gray-200 dark:text-white dark:hover:bg-slate-800">Log In</Link>
                <Link to='/signup' className="px-2 py-3 text-white bg-black text-white hover:text-slate-200 dark:bg-white dark:text-black">Get Started</Link>
                <div className="hover:text-gray-600">
                    <Link to='/settings'>
                        <p className="dark:text-white dark:hover:text-gray-300">
                            <FontAwesomeIcon icon={faGear} size="lg"/>
                        </p>
                    </Link>
                </div>
            </div>
        {/* </div> */}
    </div>
  );
}

export default Navbar;
