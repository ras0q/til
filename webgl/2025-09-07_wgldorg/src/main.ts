import fragment from "./glsl/main.frag"
import vertex from "./glsl/main.vert"
import { matIV } from "./minMatrix"

const c = document.getElementById("canvas") as HTMLCanvasElement

c.width = 300
c.height = 300

const gl = c.getContext("webgl")!

gl.clearColor(0.0, 0.0, 0.0, 1.0)
gl.clearDepth(1.0)
gl.clear(gl.COLOR_BUFFER_BIT)

const fragmentShader = gl.createShader(gl.FRAGMENT_SHADER)!
gl.shaderSource(fragmentShader, fragment)
gl.compileShader(fragmentShader)
if (!gl.getShaderParameter(fragmentShader, gl.COMPILE_STATUS)) {
  throw "fragmentShader compile error"
}

const vertexShader = gl.createShader(gl.VERTEX_SHADER)!
gl.shaderSource(vertexShader, vertex)
gl.compileShader(vertexShader)
if (!gl.getShaderParameter(vertexShader, gl.COMPILE_STATUS)) {
  throw "vertexShader compile error"
}

const program = gl.createProgram()
gl.attachShader(program, vertexShader)
gl.attachShader(program, fragmentShader)
gl.linkProgram(program)
if (gl.getProgramParameter(program, gl.LINK_STATUS)) {
  gl.useProgram(program)
} else {
  alert(gl.getProgramInfoLog(program))
}

const vertexPosition = [
  0.0, 1.0, 0.0,
  1.0, 0.0, 0.0,
  -1.0, 0.0, 0.0
]
const vbo = gl.createBuffer()
gl.bindBuffer(gl.ARRAY_BUFFER, vbo)
gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(vertexPosition), gl.STATIC_DRAW)

gl.bindBuffer(gl.ARRAY_BUFFER, vbo)

const attLocation = gl.getAttribLocation(program, "position")
const attStride = 3 // xyz
gl.enableVertexAttribArray(attLocation);
gl.vertexAttribPointer(attLocation, attStride, gl.FLOAT, false, 0, 0);

// matIVオブジェクトを生成
var m = new matIV();

// 各種行列の生成と初期化
var mMatrix = m.identity(m.create());
var vMatrix = m.identity(m.create());
var pMatrix = m.identity(m.create());
var mvpMatrix = m.identity(m.create());

// ビュー座標変換行列
m.lookAt([0.0, 1.0, 3.0], [0, 0, 0], [0, 1, 0], vMatrix);

// プロジェクション座標変換行列
m.perspective(90, c.width / c.height, 0.1, 100, pMatrix);

// 各行列を掛け合わせ座標変換行列を完成させる
m.multiply(pMatrix, vMatrix, mvpMatrix);
m.multiply(mvpMatrix, mMatrix, mvpMatrix);

// uniformLocationの取得
var uniLocation = gl.getUniformLocation(program, 'mvpMatrix');

// uniformLocationへ座標変換行列を登録
gl.uniformMatrix4fv(uniLocation, false, mvpMatrix);

gl.drawArrays(gl.TRIANGLES, 0, 3);
gl.flush()
