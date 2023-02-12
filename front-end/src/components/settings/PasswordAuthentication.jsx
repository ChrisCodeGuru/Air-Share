import { Fragment, useState } from "react";
import axios from "axios";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faXmark } from '@fortawesome/free-solid-svg-icons'
import QRCode from "react-qr-code";
import useSWR from "swr";

function PasswordAuthentication() {
    const [OPassword, setOPassword] = useState("")
    const [NPassword, setNPassword] = useState("")
    const [CPassword, setCPassword] = useState("")
    const [error, setError] = useState("")

    const [setup2fa, setSetup2fa] = useState(null)
    const [code2fa, setCode2fa] = useState("")
    const [delete2fa, setDelete2fa] = useState(false)

    const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]

	const fetcher = (url, token) =>
    axios
        .get(url, { headers: { "X-CSRF-TOKEN": csrfToken } })
        .then((res) => res.data);

    const { data: content, mutate: contentMutate } = useSWR(`/api/authentication`, fetcher);

    const handleGenerate2fa = async() => {
        try {
            const response = await axios({
                method: "get",
                url: "/api/otp/generate",
                withCredentials: true
            })

            if (response.status === 200) {
                setSetup2fa(response.data);
            }
        } catch (error) {
            if (error.response.status === 500) {
                alert("something went wrong");
            } else {
                alert(error.response.data)
            };
        }
    }

    const handleDelete2faSwitch = async() => {
        setDelete2fa(!delete2fa)
    }

    const handleVerify2faForm = async() => {
        try {
            const response = await axios({
                method: "post",
                url: "/api/otp/verify",
                data: {
                    Token: code2fa
                },
                withCredentials: true
            })

            if (response.status === 201) {
                alert("successfully setup")
                setSetup2fa(null);
                setCode2fa("");
                contentMutate();
            }
        } catch (error) {
            if (error.response.status === 500) {
                alert("something went wrong");
            } else {
                alert(error.response.data)
            };
        }
    }

    const handleDelete2faForm = async() => {
        try {
            const response = await axios({
                method: "post",
                url: "/api/otp/delete",
                data: {
                    Token: code2fa
                },
                withCredentials: true
            })

            if (response.status === 201) {
                alert("successfully deleted")
                handleDelete2faSwitch();
                setCode2fa("");
                contentMutate();
            }
        } catch (error) {
            if (error.response.status === 500) {
                alert("something went wrong");
            } else {
                alert(error.response.data)
            };
        }
    }

    async function changePasswordHandler() {
        const csrfToken = document.cookie.split('; ').find((row) => row.startsWith('csrf='))?.split('=')[1]
        if (OPassword === "" || NPassword === "" || CPassword === "") {
            setError("Please fill in all the information");
        } else {
            axios.defaults.headers.post["X-CSRF-TOKEN"] = csrfToken
            axios.post(
                "https://localhost/api/changePassword",
                JSON.stringify({
                    opassword: OPassword,
                    npassword: NPassword,
                    cpassword: CPassword,
                }),
                {withCredentials: true}
            )
            .then((res) => {
                if (res.status === 401) {
                    window.location.replace("/login");
                } else if (res.status === 200) {
                    setError("")
                    alert("Password Changed Successfully")
                } else {
                    setError(res.data)
                }
                setOPassword("")
                setNPassword("")
                setCPassword("")
            })
            .catch((err) => setError(err.response.data));
        }
    }

    return (
        <Fragment>
            <div className="p-4 w-full">
                <div className="mb-16">
                    <p className="w-full text-3xl border-b py-2 mb-4">Change Password</p>
                    {content &&
                        <form onSubmit={(e) => {e.preventDefault(); changePasswordHandler()}}>
                            {content.Password &&
                                <Fragment>
                                    <p className="mb-1">Old Password</p>
                                    <input type='password' className='border border-2 border-black py-0.5 px-2 rounded-md mb-4 text-gray-700 mb-3 w-[300px] bg-gray-200' value={OPassword} onChange={(e) => setOPassword(e.target.value)}></input>
                                </Fragment>
                            }
                            <p>New Password</p>
                            <input type='password' className='border border-2 border-black py-0.5 px-2 rounded-md mb-4 text-gray-700 mb-3 w-[300px] bg-gray-200' value={NPassword} onChange={(e) => setNPassword(e.target.value)}></input>
                            <p>Confirm Password</p>
                            <input type='password' className='border border-2 border-black py-0.5 px-2 rounded-md mb-4 text-gray-700 mb-2 w-[300px] bg-gray-200' value={CPassword} onChange={(e) => setCPassword(e.target.value)}></input>
                            <br/>
                            <div className="text-red-500 text-sm mb-3 px-2">{error}</div>
                            <button type="submit" className="border border-2 border-black py-0.5 px-3 rounded-md hover:bg-gray-100 dark:bg-gray-300 dark:text-black dark:hover:bg-gray-400">Change Password</button>
                        </form>
                    }
                </div>
                <div className="flex flex-col justify-start items-start w-full">
                    <p className="w-full text-3xl border-b py-2 mb-4">Two Factor Authentication</p>
                    {(content && !content.OTP) &&
                        <div 
                            className="cursor-pointer bg-blue-500 rounded-xl px-3 py-2 text-white font-semibold"
                            onClick={() => handleGenerate2fa()}
                        >
                            Generate 2FA
                        </div>
                    }
                    {(content && content.OTP) &&
                        <div 
                            className="cursor-pointer bg-red-500 rounded-xl px-3 py-2 text-white font-semibold"
                            onClick={() => handleDelete2faSwitch()}
                        >
                            Delete 2FA
                        </div>
                    }
                </div>
            </div>

            {/* Create form modal */}
                {setup2fa !== null &&
                    <div className="fixed w-screen h-screen flex flex-col justify-center items-center z-10 bg-slate-900/50 inset-0">
                        <div className="flex flex-col justify-center items-start bg-white text-gray-900 dark:bg-gray-800 dark:text-white rounded-2xl p-3">
                            <div className="flex flex-row justify-between items-start w-full mb-8">
                                <p className="text-lg font-semibold pl-1">Setup 2FA</p>
                                <FontAwesomeIcon 
                                    icon={faXmark} 
                                    size="lg" 
                                    className="cursor-pointer pt-0.5 pr-1.5" 
                                    onClick={() => setSetup2fa(null)}
                                />
                            </div>
                            <div className="flex justify-center items-center w-full mb-8">
                                <QRCode value={setup2fa.Otpauth_url} />
                            </div>
                            <input 
                                type="text"
                                id="2fa"
                                name="2fa"
                                className="text-base rounded-xl py-2 px-3 w-96 bg-gray-300 text-gray-600 dark:bg-gray-400 dark:text-gray-300 focus:bg-transparent transition duration-150 ease-in-out"
                                placeholder="6 digit code"
                                value={code2fa}
                                onChange={e => setCode2fa(e.target.value)}
                            />
                            <button 
                                className="w-full bg-blue-500 rounded-xl py-2 mt-2"
                                onClick={() => handleVerify2faForm()}
                            >
                                <p className="cursor-pointer text-base text-white font-semibold">
                                    Verify
                                </p>
                            </button>
                        </div>
                    </div>
                }

                {delete2fa === true &&
                    <div className="fixed w-screen h-screen flex flex-col justify-center items-center z-10 bg-slate-900/50 inset-0">
                        <div className="flex flex-col justify-center items-start bg-white text-gray-900 dark:bg-gray-800 dark:text-white rounded-2xl p-3">
                            <div className="flex flex-row justify-between items-start w-full mb-3">
                                <p className="text-lg font-semibold pl-1">Delete 2FA</p>
                                <FontAwesomeIcon 
                                    icon={faXmark} 
                                    size="lg" 
                                    className="cursor-pointer pt-0.5 pr-1.5" 
                                    onClick={() => handleDelete2faSwitch()}
                                />
                            </div>
                            <input 
                                type="text"
                                id="2fa"
                                name="2fa"
                                className="text-base rounded-xl py-2 px-3 w-96 bg-gray-300 text-gray-600 dark:bg-gray-400 dark:text-gray-300 focus:bg-transparent transition duration-150 ease-in-out"
                                placeholder="6 digit code"
                                value={code2fa}
                                onChange={e => setCode2fa(e.target.value)}
                            />
                            <button 
                                className="w-full bg-red-500 rounded-xl py-2 mt-2"
                                onClick={() => handleDelete2faForm()}
                            >
                                <p className="cursor-pointer text-base text-white font-semibold">
                                    Delete
                                </p>
                            </button>
                        </div>
                    </div>
                }
            </Fragment>
    );
}

PasswordAuthentication.defaultProps = {
    twoFactor: true
}

export default PasswordAuthentication;