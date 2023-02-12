import { Fragment, useState, useEffect } from "react";
import Navbar from "../components/navbar";
import { Link } from "react-router-dom";
import axios from "axios";

function Login() {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [error, setError] = useState("");
	const [stage, setStage] = useState(0);
	const [mfa, setMfa] = useState("");

	document.title = "Login";

	if (localStorage.theme === "dark") {
		document.documentElement.classList.add("dark");
	} else {
		document.documentElement.classList.remove("dark");
	}

	useEffect (() => {
		axios.get('/api/doubleSubmit')
		.then((res) => localStorage.setItem('csrf', res.data))
		.catch(err => localStorage.setItem('csrf', ''))
	},[]);

	async function loginHandler() {
		if (username === "" || password === "") {
			setError("Please fill in all the information");
		} else {
			axios.defaults.headers.post["X-CSRF-TOKEN"] = localStorage.csrf
			axios.post(
				"/api/login",
				JSON.stringify({
				username: username,
				password: password,
				}),
				{withCredentials: true}
			)
			.then((res) => {
				if (res.status === 401) {
					setError(res.data);
				} else if (res.status === 500){
					setError("something went wrong");
				} else if (res.status === 206){
					setStage(1);
				} else if (res.status === 200){
					window.location.replace("/fileshare");
				}
			})
			.catch((err) => setError(err.response.data));
		}
	}

	async function loginMFAHandler() {
		if (username === "" || password === "") {
			setError("Please fill in all the information");
		} else {
			axios.defaults.headers.post["X-CSRF-TOKEN"] = localStorage.csrf
			axios.post(
				"/api/login/mfa",
				JSON.stringify({
				username: username,
				password: password,
				token: mfa,
				}),
				{withCredentials: true}
			)
			.then((res) => {
				if (res.status === 401) {
					setError(res.data);
				} else if (res.status === 500){
					setError("something went wrong");
				} else if (res.status === 200){
					window.location.replace("/fileshare");
				}
			})
			.catch((err) => setError(err.response.data));
		}
	}

	return (
		<Fragment>
			<Navbar />
			<div className="bg-gray-200 min-h-screen min-w-screen flex items-center justify-center  dark:bg-gray-800">
				<div className="bg-white p-8 w-[500px] border-2 border-black shadow-xl">
					{stage === 0 &&
						<form
							onSubmit={(e) => {
							e.preventDefault();
							loginHandler();
							}}
						>
							<p className="text-2xl font-bold">Sign In</p>
							<p className="pt-2 pb-4 font-bold">
								New User?{" "}
								<Link to="/signup" className="text-blue-400 hover:text-blue-500">
									Create an account
								</Link>
							</p>
							<input
								type="text"
								placeholder="Username or Email"
								className="border-2 border-black w-full py-2 px-3 mb-4 text-gray-700 focus:outline-none focus:border-2"
								value={username}
								onChange={(e) => setUsername(e.target.value)}
							/>
							<br />
							<input
								type="password"
								placeholder="Password"
								className="border-2 border-black w-full py-2 px-3 mb-4 text-gray-700 focus:outline-none focus:border-2"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
							/>
							<br />
							<div className="text-red-500 text-sm mb-5 px-2">{error}</div>
							<button
								type="submit"
								className="bg-black w-full py-2 px-3 text-white mb-3 focus:outline-none focus:border-2"
							>
								Login
							</button>
							<a 
								className="cursor-pointer select-none flex justify-center items-center py-2 px-3 bg-blue-500 text-white"
								href="/api/google-login"
							>
								Sign in with Google
							</a>
						</form>
					}
					{stage === 1 &&
						<form
							onSubmit={(e) => {
							e.preventDefault();
							loginMFAHandler();
							}}
						>
							<p className="text-2xl font-bold pb-8">2FA</p>
							<input
								type="text"
								placeholder="6 digit code"
								className="border-2 border-black w-full py-2 px-3 mb-4 text-gray-700 focus:outline-none focus:border-2"
								value={mfa}
								onChange={(e) => setMfa(e.target.value)}
							/>
							<br />
							<div className="text-red-500 text-sm mb-5 px-2">{error}</div>
							<button
								type="submit"
								className="bg-black w-full py-2 px-3 text-white mb-3 focus:outline-none focus:border-2"
							>
								Login
							</button>
						</form>
					}
				</div>
			</div>
			<div className="px-4 pt-2 dark:bg-gray-900">
				<div className="border-t border-black flex items-center px-2 pb-2 text-xs text-2xl dark:border-gray-900 dark:text-slate-300 ">
					<div>Copyright &copy; 2022 by AirShare. All Rights Reserved</div>
				</div>
			</div>
		</Fragment>
	);
}

export default Login;
