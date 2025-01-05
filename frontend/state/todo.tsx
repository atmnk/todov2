import { createTodo, deleteTodo } from "@/services/todo";
import { atom } from "jotai";

export const todosAtom = atom<Record<number,any>>({})
export const setBulkTodosAtom = atom(null, (get, set, payload: {todos:any[]}) => {
    const todosVal:Record<number,any>={}
    console.log("Todos",payload.todos)
    for (const todo of payload.todos){
        todosVal[todo.id]=todo
    }
    console.log("Val",todosVal)
    set(todosAtom, todosVal)
})
export const createTodoAtom = atom(null, async (get, set, payload: { title: string }) => {
    const todo = await createTodo(payload.title)
    const todosVal={...get(todosAtom)}
    todosVal[todo.id]=todo
    set(todosAtom, todosVal)
    return todo

})
export const deleteTodoAtom = atom(null, async (get, set, payload: { id: number }) => {
    const todo = await deleteTodo(payload.id)
    const todosVal={...get(todosAtom)}
    delete todosVal[payload.id]
    set(todosAtom, todosVal)
    return todo

})