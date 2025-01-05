'use client'
import { deleteTodo } from "@/services/todo";
import { Button } from "@nextui-org/button";
import { Listbox, ListboxItem } from "@nextui-org/listbox";
import { redirect } from 'next/navigation'
export default function List({ todos }: { todos: { id: number, title: string, status: string }[] }) {
    return <Listbox>
        {todos.map(todo => <ListboxItem key={todo.id} id={""+todo.id}>
            <SingleTodo {...todo} />
        </ListboxItem>)}
    </Listbox>
}
function SingleTodo({ title, status, id }: { title: string, status: string, id: number }) {
    return <div className="flex flex-row items-center justify-left">
        <div className="grow">{title}</div>
        <div><Button onClick={() => {
            deleteTodo(id).then(() => {
                redirect('/')
            })
        }} color="danger">Delete</Button></div>
    </div>
}