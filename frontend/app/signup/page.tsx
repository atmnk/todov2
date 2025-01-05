import { Button } from "@nextui-org/button";
import { Card, CardBody, CardFooter, CardHeader } from "@nextui-org/card";
import { Input } from "@nextui-org/input";
import Link from "next/link";

export default function SignUp(){
    return <Card className="w-[400px] mx-auto">
        <CardHeader>Sign Up</CardHeader>
        <CardBody className="flex flex-col gap-2">
            <Input placeholder="Username"/>
            <Input placeholder="Password" type="password"/>
            <Input placeholder="Confirm Password" type="password"/>
        </CardBody>
        <CardFooter className="flex flex-row items-center justify-right gap-2">
            <Link href="/login"><Button>Login</Button></Link>
            <Button className="grow" color="primary">Sign Up</Button>
        </CardFooter>
    </Card>
}