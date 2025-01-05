import Create from "@/components/todo/create";
import List from "@/components/todo/list";
import { getAllTodos } from "@/services/todo";
import { cookies } from "next/headers";

export default async function Todo(){
    const store=await cookies()
    const todos = await getAllTodos(store.get('token')?.value!!)
    return (
      <div>
        <Create/>
        <List todos={todos}/>
      </div>
    );
}