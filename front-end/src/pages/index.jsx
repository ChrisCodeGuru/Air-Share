import { Fragment, useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Navbar from "../components/navbar";
import AuthenticatedNavbar from "../components/authenticated-navbar";
import axios from 'axios'


function Index() {
  const [display, setDisplay] = useState(1);
  const [username, setUsername] = useState('')

  document.title = "Home Page"

  useEffect (() => {
    const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]
    axios.defaults.headers.get["X-CSRF-TOKEN"] = csrfToken
    axios.get('/api/username')
    .then((res) => setUsername(res.data))
    .catch(err=> setUsername('')) 
},[]);

  if (localStorage.theme === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  } 
  
  return (
    <Fragment>
      {/* <Navbar theme={localStorage.theme}/> */}
      {username === '' ?
            <Navbar/>
        :
            <AuthenticatedNavbar username={username}/>
        }
      <div className="pb-40 dark:bg-gray-800 dark:text-white">
        <div className="text-center pt-40">
          <h1 className="text-6xl italic font-bold">File Sharing Service</h1>
          <div className="mt-10">
            <h2 className="text-3xl">
              Whether you’re sending big files for fun or delivering work for
              business
            </h2>
            <h2 className="text-3xl">projects moving forward with AirShare</h2>
          </div>

          <div className="flex flex-row justify-center items-center mt-10 dark:text-slate-300">
            <div className="cursor-pointer group flex flex-col justify-center items-center h-9" onClick={() => {setDisplay(1)}}>
              <div className="h-8 px-5">
                <p className="text-xl font-medium">Features</p>
              </div>
              <div className="bg-white group-hover:bg-black h-1 w-full rounded-sm transition ease-in-out dark:bg-gray-800 dark:group-hover:bg-white" />
            </div>
            <div className="cursor-pointer group flex flex-col justify-center items-center h-9" onClick={() => {setDisplay(2)}}>
              <div className="h-8 px-5">
                <p className="text-xl font-medium">Plans</p>
              </div>
              <div className="bg-white group-hover:bg-black h-1 w-full rounded-sm transition ease-in-out dark:bg-gray-800 dark:group-hover:bg-white" />
            </div>
            <div className="cursor-pointer group flex flex-col justify-center items-center h-9" onClick={() => {setDisplay(3)}}>
              <div className="h-8 px-5">
                <p className="text-xl font-medium">Enterprise</p>
              </div>
              <div className="bg-white group-hover:bg-black h-1 w-full rounded-sm transition ease-in-out dark:bg-gray-800 dark:group-hover:bg-white" />
            </div>
          </div>

          {display === 1 &&
            <div className="w-5/6 ml-auto mr-auto">
            {/* 1st row images */}
            <div className="container mx-auto space-y-2 lg:space-y-0 lg:gap-2 lg:grid lg:grid-cols-3 mt-10 w-full">
              <div className="w-full rounded hover:shadow-2xl">
                <img
                  className="w-screen	h-full"
                  src="https://images.unsplash.com/photo-1504711331083-9c895941bf81?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8NXx8RmlsZXxlbnwwfDB8MHx8&auto=format&fit=crop&w=600&q=60"
                  alt="pic"
                />
              </div>

              <div className="w-full rounded hover:shadow-2xl">
                <img
                  className="w-screen	h-full"
                  src="https://images.unsplash.com/photo-1481358758723-4601c3107e65?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MTV8fEZpbGV8ZW58MHwwfDB8fA%3D%3D&auto=format&fit=crop&w=600&q=60"
                  alt="pic"
                />
              </div>
              <div className="w-full rounded hover:shadow-2xl">
                <img
                  className="w-screen	h-full"
                  src="https://media.istockphoto.com/id/1347656538/photo/businesswoman-working-on-laptop-with-virtual-screen-process-automation-to-efficiently-manage.jpg?b=1&s=170667a&w=0&k=20&c=O5GVQPOsU9bWYXAVIpNqmUpqjzJI0ryja-g-EpFigDY="
                  alt="pic"
                />
              </div>
            </div>
            {/* End of 1st row images */}

            {/* 2nd row Images */}
            <div className="container mx-auto space-y-2 lg:space-y-0 lg:gap-2 lg:grid lg:grid-cols-3 mt-10">
              <div className="w-full rounded hover:shadow-2xl">
                <img className="w-screen	h-full" src="https://images.unsplash.com/photo-1586892478407-7f54fcec5b89?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MTB8fG9yZ2FuaXplfGVufDB8MHwwfHw%3D&auto=format&fit=crop&w=600&q=60" alt="pic" />
              </div>

              <div className="w-full rounded hover:shadow-2xl">
                <img className="w-screen	h-full" src="https://images.unsplash.com/photo-1630239261183-0429e79cf064?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8Zm9sZGVyfGVufDB8MHwwfHw%3D&auto=format&fit=crop&w=600&q=60" alt="pic" />
              </div>
              <div className="w-full rounded hover:shadow-2xl">
                <img className="w-screen	h-full" src="https://images.unsplash.com/photo-1457694587812-e8bf29a43845?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MTJ8fGZvbGRlcnxlbnwwfDB8MHx8&auto=format&fit=crop&w=600&q=60" alt="pic" />
              </div>
            </div>
            {/* End of 2nd row images */}
            </div>
          }

          {display === 2 &&
            <div className="w-5/6 ml-auto mr-auto flex flex-row justify-center items-center mt-10">
              <div className="flex flex-col justify-center items-center w-80 border-slate-800 border-2 rounded-xl py-5 mx-5 dark:border-slate-500">
                <h1 className="text-3xl font-bold mb-4">Standard</h1>
                <ul className="text-lg">
                  <li>Feature 1</li>
                  <li>Feature 2</li>
                  <li>Feature 3</li>
                  <li>Feature 4</li>
                  <li>Feature 5</li>
                </ul>
                <h1 className="text-2xl font-bold m-4 mb-6">Free</h1>
                <Link to='/signup' className="px-2 py-3 text-white bg-black hover:text-slate-200 dark:bg-white dark:text-black dark:hover:bg-slate-200">Get Started</Link>
              </div>
              <div className="flex flex-col justify-center items-center w-80 border-slate-800 border-2 rounded-xl py-5 mx-5 dark:border-slate-500">
                <h1 className="text-3xl font-bold mb-4">Pro</h1>
                <ul className="text-lg">
                  <li>Feature 1</li>
                  <li>Feature 2</li>
                  <li>Feature 3</li>
                  <li>Feature 4</li>
                  <li>Feature 5</li>
                </ul>
                <h1 className="text-2xl font-bold m-4 mb-6">$15<span className="text-sm font-normal">/month</span></h1>
                <Link to='/signup' className="px-2 py-3 text-white bg-black hover:text-slate-200 dark:bg-white dark:text-black dark:hover:bg-slate-200">Get Started</Link>
              </div>
              <div className="flex flex-col justify-center items-center w-80 border-slate-800 border-2 rounded-xl py-5 mx-5 dark:border-slate-500">
                <h1 className="text-3xl font-bold mb-4">Enterprise</h1>
                <ul className="text-lg">
                  <li>Feature 1</li>
                  <li>Feature 2</li>
                  <li>Feature 3</li>
                  <li>Feature 4</li>
                  <li>Feature 5</li>
                </ul>
                <h1 className="text-2xl font-bold mt-4 mb-6">Custom Pricing</h1>
                <div className="cursor-pointer w-32 px-2 py-3 text-white bg-black hover:text-slate-200 dark:bg-white dark:text-black dark:hover:bg-slate-200" onClick={() => {setDisplay(3)}}>Find out more</div>
              </div>
            </div>
          }

          {display === 3 &&
            <div className="w-5/6 ml-auto mr-auto flex flex-col justify-center items-center mt-10">
              
            </div>
          }
        </div>
      </div>
    </Fragment>
  );
}

export default Index;
