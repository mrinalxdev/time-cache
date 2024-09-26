import Navbar from "./components/Navbar";

function App() {
  return (
    <>
      <Navbar />

      <section className="mt-[4rem] text-center">
        <h1 className="text-[2rem] lg:text-[4rem] font-bold font-mono">
          TimeCache
        </h1>
        <p>
          Efficient, Thread-Safe caching with TTL for{" "}
          <span className="text-blue-500 font-bold">Golang</span> Applications
        </p>

        <p className="mt-7 max-w-3xl mx-auto">
          TimeCache is a powerful, flexible caching solution designed to engance
          the performance of your Go applications. With built-in
          time-to-live(TTL) functionality and thread-safe operations, TimeCache
          is the perfect tool for managing temporary data and reducing load on
          your primary data stores
        </p>
      </section>

      <section className="max-w-6xl mx-auto mt-[2rem] px-3">1. Installation</section>
    </>
  );
}

export default App;
