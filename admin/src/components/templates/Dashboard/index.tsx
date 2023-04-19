import React, {FC, PropsWithChildren, useState} from 'react'
import {Layout, theme} from "antd";

import './styled.sass'
import {Sidebar} from "@components/organisms";
import {WindowLocation} from "@reach/router";

type DashboardLayoutProps = {
    header?: React.ReactNode
    location: WindowLocation<WindowLocation["state"]>
}

type Props = PropsWithChildren<DashboardLayoutProps>

const DashboardLayout: FC<Props> = (props) => {
    const {children, header, location} = props

    const {
        token: { colorBgContainer },
    } = theme.useToken();

    const [collapsed, setCollapsed] = useState(false);

    return (
        <Layout className="layout">
            <Layout.Sider trigger={<>{collapsed ? '>' : '<'}</>} collapsible collapsed={collapsed} onCollapse={(value) => setCollapsed(value)}>
                <Sidebar location={location} />
            </Layout.Sider>
            <Layout>
                {header && (
                    <Layout.Header style={{background: colorBgContainer}} className="header">
                        {header}
                    </Layout.Header>
                )}
                <Layout.Content>
                    {children}
                </Layout.Content>
            </Layout>
        </Layout>
    )
}

export default DashboardLayout
