'use client'

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button, ButtonText, Field } from "@/components/component"

export default function Home() {
    const [isLogin, setLogin] = useState(true)
    const [showVerify, setVerify] = useState(false)
    const [email, setEmail] = useState("")

    return (
        <main className="min-h-screen flex flex-col items-center justify-center">
            <h1 className="text-4xl font-bold mb-10">Authentication demo page</h1>
            {
                showVerify ? <Verify email={email}/> : isLogin 
                    ? <Login setLogin={setLogin} setVerify={setVerify}/>
                    : <Register setLogin={setLogin} setVerify={setVerify} email={email} setEmail={setEmail}/>
            }
        </main>
    )
}

function Register({ setLogin, setVerify, email, setEmail }:
    { setLogin: Function, setVerify: Function, email: string, setEmail: Function }) {

    const [name, setName] = useState("")
    const [password, setPassword] = useState("")
    const [repassword, setRepassword] = useState("")
    const [error, setError] = useState("")

    const register = () => {
        if (password !== repassword) {
            return
        }

        const user = {
            name,
            email,
            password
        }

        fetch("http://localhost:3001/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(user)
        })
            .then((resp) => resp.json())
            .then(({ success, error }) => {
                if (success) {
                    setVerify(true);
                } else {
                    // handle incorrect login 
                    setError(error)
                }
            })
    }

    return (
        <>
            <h1 className="text-4xl font-bold mb-10">Register</h1>
            <div className="flex flex-col gap-5">
                <Field title={'Name'} value={name} onChange={setName} />
                <Field value={email} onChange={setEmail}
                    title={'Email'} placeholder="example@mail.com" type="email" />
                <Field value={password} onChange={setPassword}
                    title={'Password'} type="password" />
                <Field value={repassword} onChange={setRepassword}
                    title={'Retype Password'} type="password" />
            </div>
            <Button text={"Register"} onClick={register} />
            <ButtonText text={"Login instead"} onClick={() => setLogin(true)} />
        </>

    )
}

function Login({ setLogin, setVerify }: { setLogin: Function , setVerify: Function }) {
    const router = useRouter()

    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState("")

    const login = () => {
        const user = {
            email,
            password
        }

        fetch("http://localhost:3001/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(user)
        })
            .then((resp) => resp.json())
            .then(({ success, isVerified, error }) => {
                if (success) {
                    if (isVerified) {
                        router.push("/dashboard")
                    }
                    setVerify(true);
                } else {
                    // handle incorrect login 
                    setError(error)
                }
            })
    }

    return (

        <>
            <h1 className="text-4xl font-bold mb-10">Login</h1>
            <div className="flex flex-col gap-5">
                <Field value={email} onChange={setEmail}
                    title={'Email'} placeholder="example@mail.com" type="email" />
                <Field value={password} onChange={setPassword}
                    title={'Password'} type="password" />
            </div>
            {error === "" ? <></> : <p className="my-2 text-red-800">{error}</p>}
            <Button text={"Login"} onClick={() => { login() }} />
            <ButtonText text={"Register instead"} onClick={() => setLogin(false)} />
        </>
    )
}

function Verify({ email }: {email: string}) {

    const router = useRouter()
    const [code, setCode] = useState("")
    const [error, setError] = useState("")

    const verify = () => {
        fetch(`http://localhost:3001/verify?email=${email}&code=${code}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
        })
            .then((resp) => resp.json())
            .then(({ success, error }) => {
                if (success) {
                    router.push("/dashboard")
                } else {
                    // handle incorrect login 
                    setError(error)
                }
            })
    }

    return (
        <div className="flex flex-col gap-5 items-center justify-center">
            <h1 className="uppercase text-xl" >Verify your email address</h1>
            <div className="w-1/2 text-center">
            <p className="">A verification code and link has been sent to *****@****mail.com</p>
            </div>

            <p>Please check your inbox or junk folder and proceed to verify with us. The code will expire in 2mins</p>

            <Field title="Code" value={code} onChange={setCode} />

            <Button text={"Verify"} onClick={() => {verify()}} />
            <Button text={"resend code"} onClick={() => {}} />
        </div>
    )
}
