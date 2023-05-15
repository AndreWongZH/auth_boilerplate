"use client"

import { Button } from "@/components/component";
import Link from "next/link";

export default function Page() {
    return (
        <div className="min-h-screen flex flex-col items-center justify-center">
            <h1 className="text-3xl mb-3">You are now verified!</h1>
            <h2>Click below to return back to the homepage and login</h2>
            <Link href="/">
                <Button text="Return"/>
            </Link>
        </div>
    )
}
