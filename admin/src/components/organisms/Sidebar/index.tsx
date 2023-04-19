import React, {FC} from 'react';
import {Menu} from "antd";
import type {MenuProps} from 'antd';
import {Link} from "gatsby";
import {WindowLocation} from "@reach/router";

type MenuItem = Required<MenuProps>['items'][number];

function getItem(
    label: React.ReactNode,
    key: React.Key,
    icon?: React.ReactNode,
    children?: MenuItem[],
    type?: 'group',
): MenuItem {
    return {
        key,
        icon,
        children,
        label,
        type,
    } as MenuItem;
}


const items: MenuItem[] = [
    getItem(<Link to="/">Inicio</Link>, '/', undefined),

    getItem('VPN', '/vpn', undefined, [
        getItem(<Link to="/vpn/users">Usuarios</Link>, '/vpn/users', undefined),
    ]),

    getItem('AWS', '/aws', undefined, [
        getItem(<Link to="/aws/ecs">ECS</Link>, '/aws/ecs', undefined),
    ]),

    { type: 'divider' },

    getItem('Configuración', '/settings', undefined, [
        getItem(<Link to="/settings/users">Usuarios</Link>, '/settings/users/', undefined),
        getItem(<Link to="/settings/modules">Modulos</Link>, '/settings/modules/', undefined),
    ]),

];

type Props = {
    location: WindowLocation<WindowLocation["state"]>
}

const Sidebar: FC<Props> = (props) => {
    const {location} = props;

    return <>
        <Menu
            selectedKeys={[location.pathname]}
            theme="dark"
            mode="inline"
            defaultSelectedKeys={['/']}
            items={items}
            defaultOpenKeys={[location.pathname.substring(0, location.pathname.substring(1).indexOf('/') + 1)]}
        />
    </>
};

export default Sidebar;
