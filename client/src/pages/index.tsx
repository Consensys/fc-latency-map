import type { NextPage } from "next";
import Head from "next/head";
import Image from "next/image";
import styles from "../styles/Home.module.css";
import dynamic from "next/dynamic";
import Map from "../components/Map";

import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";

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
        <meta name="description" content="Filecoin Latency Map" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className={styles.main} style={{ width: "100%", height: "100%" }}>
        <div id="map" style={{ width: "100%", height: "100%" }}>
          <MapWithNoSSR data={data} />
        </div>
      </main>

      <footer className={styles.footer}>
        <a
          href="https://vercel.com?utm_source=create-next-app&utm_medium=default-template&utm_campaign=create-next-app"
          target="_blank"
          rel="noopener noreferrer"
        >
          Filecoin Latency Map
        </a>
      </footer>
    </div>
  );
};

export async function getServerSideProps(context) {
  const data = await import("../../data/export.json");

  // console.log(data1);

  // const data = [
  //   [40.8054, -74.0241],
  //   [52.48395, -1.8898],
  // ];

  return {
    props: { data: JSON.stringify(data) }, // will be passed to the page component as props
  };
}

export default Home;
