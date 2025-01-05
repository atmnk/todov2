import { getBaseURL } from "@/lib/helper";

export async function createTodo(title:string){
    const resp=await fetch(`${getBaseURL()}/todo`,{
        method:'POST',
        body:JSON.stringify({title})
    })
    if (resp.status==201){
        return (await resp.json()).data
    } else {
        const body =await resp.json()
        throw Error(body.message+" "+body.data)
    }
}
export async function deleteTodo(id:number){
    return fetch(`${getBaseURL()}/todo/${id}`,{
        method:'DELETE',
    })
}
export async function getAllTodos(token:string){
    const resp=await fetch(`${getBaseURL()}/todo`,{
        headers:{
            'Authorization':`Bearer ${token}`
        }
    })
    const data=await resp.json()
    return data.data;
}