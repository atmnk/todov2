import Create from "@/components/todo/create";
import List from "@/components/todo/list";
import { getAllTodos } from "@/services/todo";

export default async function Home() {
  const todos = await getAllTodos()
  return (
    <div>
      <Create/>
      <List todos={todos}/>
    </div>
  );
}
