import React from 'react';
import {Button, Form, Input, Row} from "antd";
import './styled.sass';

type FormLoginProps = {
    username: string;
    password: string;
};

const initialValues: FormLoginProps = {
    username: '',
    password: '',
};

const FormLogin = () => {
    const [form] = Form.useForm();

    const handleSubmit = (values: FormLoginProps) => {
        console.log(values);
    }

    return (
        <Row justify="center">
            <Form
                className="form-login"
                form={form}
                initialValues={initialValues}
                onFinish={handleSubmit}
                layout="vertical">
                <Form.Item label="Usuario" name="username">
                    <Input/>
                </Form.Item>
                <Form.Item label="Contraseña" name="password">
                    <Input.Password/>
                </Form.Item>
                <Form.Item>
                    <Button type="primary" htmlType="submit">
                        Iniciar sesión
                    </Button>
                </Form.Item>
            </Form>
        </Row>
    );
};

export default FormLogin;
