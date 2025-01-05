'use client'
import { createTodo } from "@/services/todo";
import { Button } from "@nextui-org/button";
import { Input } from "@nextui-org/input";
import { redirect } from "next/navigation";
import { useState } from "react";

export default function Create() {
    const [title,setTitle] = useState('')
    return <div className="flex flex-row gap-2">
        <Input placeholder="Todo Title" className="grow" value={title} onChange={(e)=>{setTitle(e.target.value)}}/>
        <Button onClick={()=>{
            createTodo(title).then(()=>{
                setTitle('')
                redirect('/')
            })
        }}>Add</Button>
    </div>
}