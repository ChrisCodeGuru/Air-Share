import { useState, useEffect } from "react";
import axios from "axios";

function Account() {
    const [fname, setFname] = useState('')
    const [lname, setLname] = useState('')
    const [username, setUsername] = useState('')
    const [updateError, setUpdateError] = useState('')
    const [deleteAccount, setDeleteAccount] = useState("")
    const [deleteError, setDeleteError] = useState("")
    const [userInfo, setUserInfo] = useState({})

    useEffect (() => {
        const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]
        axios.defaults.headers.get["X-CSRF-TOKEN"] = csrfToken
        axios.get('/api/userinfo')
        .then((res) => {
            if (res.status === 200){
                setUserInfo(res.data)
            }
        })
        .catch(err=> setUserInfo(''))
    },[]);

    async function updateAccountHandler() {
        const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]
        axios.defaults.headers.post["X-CSRF-TOKEN"] = csrfToken
        axios.post(
            "/api/updateAccount",
            JSON.stringify({
                fname: fname,
                lname: lname,
                username: username,
                email: userInfo.email
            }),
            {withCredentials: true}
        )
        .then((res) => {
            if (res.status === 401) {
                window.location.replace("/login");
            } else if (res.status === 200) {
                alert("Account Updated Successfully")
            } else {
                setUpdateError(res.data)
            }
        })
        .catch((err) => setUpdateError(err.response.data));
      }

    async function deleteAccountHandler() {
        const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]
        if (deleteAccount === "") {
            setDeleteError("Please fill in all the information");
        } else {
            axios.defaults.headers.post["X-CSRF-TOKEN"] = csrfToken
            axios.post(
                "/api/deleteAccount",
                JSON.stringify({
                    delete: deleteAccount
                }),
                {withCredentials: true}
            )
            .then((res) => {
                if (res.status === 401) {
                    window.location.replace("/login");
                } else if (res.status === 200) {
                    alert("Account Deleted Successfully")
                    window.location.replace("/");
                } else {
                    setDeleteError(res.data)
                }
            })
            .catch((err) => setDeleteError(err.response.data));
        }
      }

    return (
        <div className="p-4 w-full">
            <div className="mb-16">
                <p className="w-full text-3xl border-b py-2 mb-4">Edit Account Information</p>
                <form onSubmit={(e) => {e.preventDefault(); updateAccountHandler()}}>
                    <p className="mb-0.5 font-semibold">First Name:</p>
                    <input placeholder={userInfo.fname}  className='border-[1px] border-black px-2 py-0.5 w-[300px] dark:text-black' value={fname} onChange={e => setFname(e.target.value)}></input>
                    <p className="mb-0.5 font-semibold">Last Name:</p>
                    <input placeholder={userInfo.lname}  className='border-[1px] border-black px-2 py-0.5 w-[300px] dark:text-black' value={lname} onChange={e => setLname(e.target.value)}></input>
                    <p className="mb-0.5 font-semibold">Username:</p>
                    <input placeholder={userInfo.username}  className='border-[1px] border-black px-2 py-0.5 w-[300px] dark:text-black' value={username} onChange={e => setUsername(e.target.value)}></input>
                    <p className="mb-0.5 font-semibold">Email:</p>
                    <p className='border-[1px] border-black px-2 py-0.5 w-[300px] dark:text-black'>{userInfo.email}</p>
                    <div className="text-red-500 text-sm pb-0.5 px-2">{updateError}</div>
                    <button type="submit" className="my-2 py-0.5 px-3 rounded-md border-2 border-green-600 bg-green-400 hover:bg-green-500">Save</button>
                </form>
            </div>
            <div>
                <p className="w-full text-3xl border-b py-2 mb-4">Delete Account</p>
                <p className="mb-1">Enter <span className="font-bold">Delete {userInfo.username}</span> to delete account</p>
                <form onSubmit={(e) => {e.preventDefault(); deleteAccountHandler()}}>
                    <input type='text' placeholder={'Delete ' + userInfo.username} className='border-[1px] border-black rounded-md px-2 py-0.5 w-[300px] dark:text-black'  value={deleteAccount} onChange={(e) => setDeleteAccount(e.target.value)}></input>
                    <div className="text-red-500 text-sm pb-0.5 px-2">{deleteError}</div>
                    <button type="submit" className="my-2 py-0.5 px-3 rounded-md border-2 border-red-600 bg-red-400 hover:bg-red-500">Delete</button>
                </form>
            </div>
        </div>
    );
}

export default Account;