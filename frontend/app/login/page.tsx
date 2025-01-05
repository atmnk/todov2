import { Button } from "@nextui-org/button";
import { Card, CardBody, CardFooter, CardHeader } from "@nextui-org/card";
import { Input } from "@nextui-org/input";

export default function Login(){
    return <Card>
        <CardHeader>Login</CardHeader>
        <CardBody>
            <Input placeholder="Username"/>
            <Input placeholder="Password" type="password"/>
        </CardBody>
        <CardFooter>
            <Button color="primary">Login</Button>
        </CardFooter>
    </Card>
}