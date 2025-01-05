'use client'
import { deleteTodoAtom, setBulkTodosAtom, todosAtom } from "@/state/todo";
import { Button } from "@nextui-org/button";
import { Listbox, ListboxItem } from "@nextui-org/listbox";
import { useAtom, useAtomValue } from "jotai";
export default function List() {
    const currentTodos = useAtomValue(todosAtom)
    console.log(JSON.stringify(currentTodos))
    return <Listbox>
        {Object.entries(currentTodos).map(([_,todo]) => <ListboxItem key={todo.id} id={""+todo.id}>
            <SingleTodo {...todo} />
        </ListboxItem>)}
    </Listbox>
}
function SingleTodo({ title, status, id }: { title: string, status: string, id: number }) {
    const [,deleteTodo] = useAtom(deleteTodoAtom)
    return <div className="flex flex-row items-center justify-left">
        <div className="grow">{title}</div>
        <div><Button onClick={() => {
            deleteTodo({id})
        }} color="danger">Delete</Button></div>
    </div>
}