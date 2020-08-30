#version 330 core

layout (location = 0) in vec2 in_position;
layout (location = 1) in vec2 in_uv;

out vec2 pass_uv;

uniform mat4 u_transform;

void main() {
	gl_Position = u_transform*vec4(in_position, 0., 1);
    gl_Position.z = clamp(gl_Position.z, -1, 1);
	pass_uv = in_uv;
}