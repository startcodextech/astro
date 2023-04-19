import React, {FC} from "react";
import {DashboardLayout} from "@components/templates";
import {PageProps} from "gatsby";

const SettingsUsersPage: FC<PageProps> = (props) => {
    const {location} = props;

    return (
        <>
            <DashboardLayout location={location}>
                <>settings</>
            </DashboardLayout>
        </>
    )
}
export default SettingsUsersPage

export const Head: FC = () => <title>Configuracion - Usuario</title>
