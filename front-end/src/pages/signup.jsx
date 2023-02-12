import { Fragment, useState } from "react";
import { Link } from "react-router-dom";
import axios from 'axios'
import Navbar from '../components/navbar'

function Signup() {
  const [fname, setFname] = useState('')
  const [lname, setLname] = useState('')
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  
  document.title = "Sign Up"

  if (localStorage.theme === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  } 

  async function createUser () {
    var re1 = /^[a-zA-Z\s]+$/
    var re2 = /^[a-zA-Z0-9 ]+$/

    if (fname === '' || lname === '' || username === '' || email === '' || password === '' || confirmPassword === '') {
      setError("Please fill in all the information")
    }
    else if (!re1.test(fname) || !re1.test(lname)) {
      setError("Name cannot contain Special Characters or Numbers")

    }
    else if(!re2.test(username)) {
      setError("Username cannot contain Special Character")
    }
    else if (password.length < 8) {
      setError("Passwords must be at least 8 Characters long")
    }
    else if (password !== confirmPassword) {
      setError("Passwords Do Not Match")
    }
    else {
      axios.post('https://localhost/api/createUser', JSON.stringify(
        {
            fname: fname,
            lname: lname,
            username: username,
            email: email,
            password: password,
            cpassword: confirmPassword
        }
      ))
      .then((res) => {
        if (res.status === 201){
          window.location.replace('/login')
        }
        else{setError(res.data)}
      })
      .catch(err=> setError(err.response.data))    
    }
  }

  return (
    <Fragment>
      <Navbar/>
      <div className="bg-gray-200 min-h-screen min-w-screen flex items-center justify-center dark:bg-gray-800">
        <div className="bg-white p-8 max-w-[500px] border-2 border-black shadow-xl">
          <form onSubmit={(e) => {e.preventDefault(); createUser()}}>
            <p className="text-2xl font-bold ">Sign Up</p>
            <p className="pt-2 pb-4 font-bold ">Have an account? <Link to='/login' className="text-blue-400 hover:text-blue-500">Log In</Link></p>
            <div className="flex justify-between">
              <input type='text' placeholder='First Name' className='border-2 border-black w-[49%] py-2 px-3 text-gray-700 mb-3 focus:outline-none focus:border-2' value={fname} onChange={e => setFname(e.target.value)}></input>
              <input type='text' placeholder='Last Name' className='border-2 border-black w-[49%] py-2 px-3 text-gray-700 mb-3 focus:outline-none focus:border-2' value={lname} onChange={e => setLname(e.target.value)}></input>
            </div>
            <input type='text' placeholder='Username' className='border-2 border-black w-full py-2 px-3 text-gray-700 mb-3 focus:outline-none focus:border-2' value={username} onChange={e => setUsername(e.target.value)}></input>
            <input type='text' placeholder='Email' className='border-2 border-black w-full py-2 px-3 text-gray-700 mb-3 focus:outline-none focus:border-2 ' value={email} onChange={e => setEmail(e.target.value)}></input>
            <input type='password' placeholder='Password' className='border-2 border-black w-full py-2 px-3 text-gray-700 mb-3 focus:outline-none focus:border-2' value={password} onChange={e => setPassword(e.target.value)}></input>
            <input type='password' placeholder='Confirm Password' className={`border-2 border-black w-full py-2 px-3 mb-3 text-gray-700 ${error === '' ? 'mb-8' : 'mb-3'} focus:outline-none focus:border-2`} value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)}></input>
            <div className="text-red-500 text-sm mb-5 px-2">{error}</div>
            <button type='submit' className='bg-black w-full py-2 px-3 text-white mb-3 focus:outline-none focus:border-2 hover:text-gray-200'>Sign Up</button>
          </form>
        </div>
      </div>
      <div className='px-4 pt-2 dark:bg-gray-900'>
        <div className='border-t border-black flex items-center px-2 pb-2 text-xs text-2xl dark:border-gray-900 dark:text-slate-300 '>
          <div>Copyright &copy; 2022 by AirShare. All Rights Reserved</div>
        </div>
      </div>
    </Fragment>
  );
}

export default Signup;
