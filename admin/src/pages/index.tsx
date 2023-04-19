import * as React from "react"
import type {HeadFC, PageProps} from "gatsby"
import {DashboardLayout} from "@components/templates";

const IndexPage: React.FC<PageProps> = (props) => {
    const {location} = props;
    return (
        <DashboardLayout location={location}>
            <>
                hola
            </>
        </DashboardLayout>
    )
}

export default IndexPage

export const Head: HeadFC = () => <title>Astro</title>
