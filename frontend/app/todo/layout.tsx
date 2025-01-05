import "@/styles/globals.css";
import { Metadata, Viewport } from "next";
import { Link } from "@nextui-org/link";
import clsx from "clsx";


import { siteConfig } from "@/config/site";
import { fontSans } from "@/config/fonts";
import { Button } from "@nextui-org/button";
import Header from "@/components/header";
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

export default function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <>
        <Header/>
        <div className="container mx-auto max-w-7xl p-6 flex-grow border my-2 rounded-xl ">
            {children}
        </div>
        </>
    );
}
