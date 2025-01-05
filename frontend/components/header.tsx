'use client'
import { Button } from "@nextui-org/button";
import { redirect } from "next/navigation";

export default function Header(){
    return <div className="flex flex-row"><div className="grow text-3xl">My Todos</div><Button onPress={()=>{
        document.cookie = "token=; path=/; secure; samesite=strict; expires=Thu, 01 Jan 1970 00:00:00 GMT";
        redirect("/")

    }}>Logout</Button></div>
}