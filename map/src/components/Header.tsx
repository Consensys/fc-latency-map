import React from "react";
import getConfig from "next/config";
import styles from "@src/styles/Global.module.css";

import { Layout } from "antd";
import Link from "next/link";
import { useRouter } from "next/router";

import { Typography } from "antd";

const { Title } = Typography;

const { publicRuntimeConfig } = getConfig();

const { Header } = Layout;

const TopHeader = () => {
  const router = useRouter();
  let selectedKeys = [router.pathname];

  return (
    <Header className={styles.header}>
      <div className={styles.headerLogo}>
        <Title level={5}>
          <Link href="/">{publicRuntimeConfig.app.name}</Link>
        </Title>
      </div>
    </Header>
  );
};

export default TopHeader;
