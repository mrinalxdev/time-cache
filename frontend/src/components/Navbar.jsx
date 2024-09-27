const Navbar = () => {
  return (
    <div className="flex justify-between items-center max-w-5xl mx-auto m-5 px-3">
      <div>
        <span className="text-xl font-bold font-mono">⏱️TimeCache</span>
      </div>

      <div className="flex gap-4">
        <a
          href="https://github.com/mrinalxdev/time-cache"
          className="hover:underline hover:underline-offset-4 duration-100 cursor-pointer"
        >
          Github
        </a>

        <a
          href="https://github.com/mrinalxdev/time-cache"
          className="hover:underline hover:underline-offset-4 duration-100 cursor-pointer"
        >
          Lists of Projects
        </a>
      </div>
    </div>
  );
};

export default Navbar;
