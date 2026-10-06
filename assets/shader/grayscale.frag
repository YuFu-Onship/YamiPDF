#version 330 core

in vec2 uv;
out vec4 FragColor;
uniform sampler2D tex;

void main()
{
    vec4 color = texture(tex, uv);
    float gray =
        0.299 * color.r +
            0.587 * color.g +
            0.114 * color.b;

    FragColor = vec4(
            gray,
            gray,
            gray,
            color.a
        );
}
