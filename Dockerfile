# FROM balenalib/rpi-raspbian
FROM balenalib/raspberrypi4-64-debian:bullseye-run-20221031

COPY bot /bot

CMD ["/bot"]