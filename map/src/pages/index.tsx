import type { NextPage } from "next";
import Head from "next/head";
import Image from "next/image";
import styles from "@src/styles/Home.module.css";
import dynamic from "next/dynamic";
import Map from "@src/components/Map";
import getConfig from "next/config";

import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";

const { publicRuntimeConfig, serverRuntimeConfig } = getConfig();

interface Props {
  data: any;
}

const Home = (props: Props) => {
  const { data } = props;
  const MapWithNoSSR = dynamic(() => import("../components/Map"), {
    ssr: false,
  });

  return (
    <div className={styles.container}>
      <Head>
        <title>Filecoin Latency Map</title>
        <meta name="description" content={publicRuntimeConfig.app.name} />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className={styles.main} style={{ width: "100%", height: "100%" }}>
        <div id="map" style={{ width: "100%", height: "100%" }}>
          <MapWithNoSSR data={data} />
        </div>
      </main>

      <footer className={styles.footer}>
        <a href="/" target="_blank" rel="noopener noreferrer">
          {publicRuntimeConfig.app.name}
        </a>
      </footer>
    </div>
  );
};

export async function getServerSideProps() {
  // const data = await import(
  //   serverRuntimeConfig.path.exportsMeasures + "export.json"
  // );
  const data = await import("../../data/data_1633076926.json");

  return {
    props: {
      data: JSON.stringify(data),
    },
  };
}

export default Home;
