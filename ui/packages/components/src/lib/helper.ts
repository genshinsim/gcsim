export function charBG(element: string) {
  switch (element) {
    case "cryo":
      return "bg-gradient-to-r from-gray-700 to-blue-300";
    case "hydro":
      return "bg-gradient-to-r from-gray-700 to-blue-500";
    case "pyro":
      return "bg-gradient-to-r from-gray-700 to-red-400";
    case "electro":
      return "bg-gradient-to-r from-gray-700 to-purple-300";
    case "anemo":
      return "bg-gradient-to-r from-gray-700 to-teal-500";
    case "dendro":
      return "bg-gradient-to-r from-gray-700 to-lime-700";
    case "geo":
      return "bg-gradient-to-r from-gray-700 to-yellow-400";
  }
  return "bg-gray-700";
}

export function prettyPrintNumberStr(num: string): string {
  return num.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}
