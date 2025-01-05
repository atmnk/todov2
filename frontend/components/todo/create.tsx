'use client'
import { createTodoAtom } from "@/state/todo";
import { Button } from "@nextui-org/button";
import { Input } from "@nextui-org/input";
import { useAtom } from "jotai";
import { redirect } from "next/navigation";
import { useState } from "react";

export default function Create() {
    const [title,setTitle] = useState('')
    const [,createTodo] = useAtom(createTodoAtom)
    return <div className="flex flex-row gap-2 p-2">
        <Input placeholder="Todo Title" className="grow" value={title} onChange={(e)=>{setTitle(e.target.value)}}/>
        <Button color="primary" onClick={()=>{
            createTodo({title}).then(()=>{
                setTitle('')
            })
        }}>Add</Button>
    </div>
}