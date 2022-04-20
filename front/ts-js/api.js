var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
export const server = {
    get: function (path) {
        return __awaiter(this, void 0, void 0, function* () {
            const port = location.port;
            console.log(port);
            const res = yield fetch("http://localhost:" + port + "/api/" + path);
            if (res.status != 200) {
                console.error(res);
                return;
            }
            return res.json();
        });
    },
};
export const sleep = (waitTime) => new Promise((resolve) => setTimeout(resolve, waitTime));
export default "default export";
