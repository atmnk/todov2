import { getBaseURL } from "@/lib/helper";

export async function signup(username:string,password:string){
    const resp=await fetch(`${getBaseURL()}/user`,{
        method:'POST',
        body:JSON.stringify({username,password})
    })
    if (resp.status==201){
        return (await resp.json()).data
    } else {
        const body=await resp.json()
        throw Error(body.message+" "+body.data)
    }
}
export async function login(username:string,password:string){
    const resp=await fetch(`${getBaseURL()}/user/login`,{
        method:'POST',
        body:JSON.stringify({username,password})
    })
    if (resp.status==200){
        return (await resp.json()).data
    } else {
        const body=await resp.json()
        throw Error(body.message+" "+body.data)
    }
}