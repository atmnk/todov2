import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export default async function Home() {
  const store = await cookies()
  const token = store.get('token')
  if (!token){
    redirect("/login")
  } else {
    redirect("/todo")
  }
}
