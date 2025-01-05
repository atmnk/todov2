import "@/styles/globals.css";
import { Metadata, Viewport } from "next";
import { redirect } from 'next/navigation';


import Header from "@/components/header";
import { siteConfig } from "@/config/site";
import { getAllTodos } from "@/services/todo";
import { cookies } from "next/headers";
export const metadata: Metadata = {
    title: {
        default: siteConfig.name,
        template: `%s - ${siteConfig.name}`,
    },
    description: siteConfig.description,
    icons: {
        icon: "/favicon.ico",
    },
};

export const viewport: Viewport = {
    themeColor: [
        { media: "(prefers-color-scheme: light)", color: "white" },
        { media: "(prefers-color-scheme: dark)", color: "black" },
    ],
};

export default async function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    const store = await cookies()
    const token = store.get('token')
    if (!token) {
        redirect("/login")
    } else {
        const todos = await getAllTodos(store.get('token')?.value!!)

        return (
            <>
                <Header todos={todos} />
                <div className="container mx-auto max-w-7xl p-6 flex-grow border my-2 rounded-xl ">
                    {children}
                </div>
            </>
        );
    }
}
