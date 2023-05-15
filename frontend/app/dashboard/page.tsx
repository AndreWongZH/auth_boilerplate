"use client"

import { Button } from "@/components/component";
import Link from "next/link";

export default function Page() {
    return (
        <div className="min-h-screen flex flex-col items-center justify-center">
            <h1 className="text-3xl mb-3">This is a dashboard viewable when logged in</h1>
            <h2>Return to root page using button below</h2>
            <Link href="/">
                <Button text="Return"/>
            </Link>
        </div>
    )
}
