import React from "react";
import {FC} from "react";
import {HeadFC, PageProps} from "gatsby";
import {FormLogin} from "@components/organisms";
import {Col, Row} from "antd";
import {Copyright, HeadLogin} from "@components/molecules";

import './styled.sass';

const LoginPage: FC<PageProps> = () => {

    return (
        <>
            <Row>
                <Col className="gutter-row" xs={0} sm={0} md={14}>
                    <div className="login-lead">
                        <div className="background"/>
                        <h1 className="lead">Tu estación espacial para manejar servicios en la nube.</h1>
                    </div>
                </Col>
                <Col className="gutter-row" xs={24} sm={24} md={10}>
                    <HeadLogin />
                    <FormLogin />
                    <Copyright />
                </Col>
            </Row>

        </>
    )

}

export default LoginPage;

export const Head: HeadFC = () => <title>Iniciar sesión</title>
