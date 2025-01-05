import { getBaseURL } from "@/lib/helper";

export async function createTodo(title:string){
    return fetch(`${getBaseURL()}/todo`,{
        method:'POST',
        body:JSON.stringify({title})
    })
}
export async function deleteTodo(id:number){
    return fetch(`${getBaseURL()}/todo/${id}`,{
        method:'DELETE',
    })
}
export async function getAllTodos(){
    const resp=await fetch(`${getBaseURL()}/todo`)
    const data=await resp.json()
    return data.data;
}