package render

// 顶点着色器
const vertexShaderSource = `
#version 330 core
layout (location = 0) in vec2 aPos;
layout (location = 1) in vec2 aTexCoord;
out vec2 TexCoord;
void main() {
    gl_Position = vec4(aPos, 0.0, 1.0);
    TexCoord = aTexCoord;
}` + "\x00"

// 黑白shader
const grayShader = `
#version 330 core
out vec4 FragColor;
in vec2 TexCoord;
uniform sampler2D uTexture;
void main() {
    vec4 texColor = texture(uTexture, TexCoord);
    // 加权灰度算法 (Luminance)
    float gray = dot(texColor.rgb, vec3(0.2126, 0.7152, 0.0722));
    FragColor = vec4(vec3(gray), texColor.a);
}` + "\x00"
