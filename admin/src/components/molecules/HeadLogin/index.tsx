import React from 'react';
import { StaticImage } from "gatsby-plugin-image";
import {Col, Row} from "antd";
import './styled.sass'

const HeadLogin = () => {
    return (
        <Row justify="center">
            <div className="login-head">
                <Row justify="center">
                    <Col>
                        <StaticImage src="../../../assets/images/favicons/favicon.png" alt="" width={50} height={50}/>
                    </Col>
                    <h1 className="title">ASTRO</h1>
                </Row>
                <div className="welcome">
                    <p>Bienvenido</p>
                    <span>Inicia sesión para acceder a tu cuenta.</span>
                </div>


            </div>
        </Row>
    );
}

export default HeadLogin;
