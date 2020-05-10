#version 330 core
in vec2 pass_uv;

out vec4 out_color;

uniform sampler2D u_tex;

void main()
{
    out_color = texture(u_tex, pass_uv);
} 