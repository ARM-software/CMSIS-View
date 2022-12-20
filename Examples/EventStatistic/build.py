#!/usr/bin/python
# -*- coding: utf-8 -*-

import logging

from datetime import datetime
from enum import Enum
from glob import glob, iglob
from pathlib import Path

from lxml.etree import XMLSyntaxError
from zipfile import ZipFile

from matrix_runner import main, matrix_axis, matrix_action, matrix_command, matrix_filter, \
    ConsoleReport, CropReport, TransformReport, JUnitReport


@matrix_axis("device", "d", "Device(s) to be considered.")
class DeviceAxis(Enum):
    CM3    = ('Cortex-M3',  'CM3')
    CM55   = ('Cortex-M55', 'CM55')
    SSE300 = ('Corstone_SSE-300', 'SSE300')


@matrix_axis("compiler", "c", "Compiler(s) to be considered.")
class CompilerAxis(Enum):
    AC6 = ('AC6')
    GCC = ('GCC')
    IAR = ('IAR')

    @property
    def image_ext(self):
        ext = {
            CompilerAxis.AC6: 'axf',
            CompilerAxis.GCC: 'elf',
            CompilerAxis.IAR: 'elf'
        }
        return ext[self]


@matrix_axis("optimize", "o", "Optimization to be considered.")
class OptimizeAxis(Enum):
    DEBUG = ('Debug')
    RELEASE  = ('Release')


MODEL_EXECUTABLE = {
    DeviceAxis.CM3: ("VHT_MPS2_Cortex-M3", []),
    DeviceAxis.CM55: ("VHT_MPS2_Cortex-M55", []),
    DeviceAxis.SSE300: ("VHT_MPS3_Corstone_SSE-300", [])
}


def config_suffix(config, timestamp=True):
    suffix = f"{config.compiler[0]}-{config.optimize}-{config.device[1]}"
    if timestamp:
        suffix += f"-{datetime.now().strftime('%Y%m%d%H%M%S')}"
    return suffix


def project_name(config):
    return f"EventStatistic.{config.compiler}_{config.optimize}+{config.device[1]}"


def output_dir(config):
    return f"{project_name(config)}_outdir"


def model_config(config):
    return f"model_config_{config.device[1].lower()}.txt"


@matrix_action
def clean(config):
    """Build the selected configurations using CMSIS-Build."""
    yield cbuild_clean(f"{project_name(config)}/{project_name(config)}.cprj")


@matrix_action
def build(config, results):
    """Build the selected configurations using CMSIS-Build."""
    logging.info("Compiling Project...")

    yield csolution(f"{project_name(config)}")
    yield cbuild(f"{project_name(config)}/{project_name(config)}.cprj")

    if not all(r.success for r in results):
        return

    file = f"blinky-{config_suffix(config)}.zip"
    logging.info(f"Archiving build output to {file}...")
    with ZipFile(file, "w") as archive:
        for content in iglob(f"{project_name(config)}/**/*", recursive=True):
            if Path(content).is_file():
                archive.write(content)


@matrix_action
def extract(config):
    """Extract the latest build archive."""
    archives = sorted(glob(f"EventStatistic-{config_suffix(config, timestamp=False)}-*.zip"), reverse=True)
    yield unzip(archives[0])


@matrix_action
def run(config, results):
    """Run the selected configurations."""
    logging.info("Running Event Statistic Example on Arm model ...")
    yield model_exec(config)


@matrix_action
def events(config, results):
    """Dump event log."""
    logging.info("Dump event log ...")
    yield eventlist(config)


@matrix_command()
def cbuild_clean(project):
    return ["cbuild", "-c", project]


@matrix_command()
def unzip(archive):
    return ["bash", "-c", f"unzip {archive}"]


@matrix_command()
def csolution(project):
    return ["csolution", "convert", "-s", "EventStatistic.csolution.yml", "-c", project]


@matrix_command()
def cbuild(project):
    return ["cbuild", project]


@matrix_command()
def model_exec(config):
    cmdline = [MODEL_EXECUTABLE[config.device][0], "-q", "--simlimit", 200, "-f", model_config(config)]
    cmdline += MODEL_EXECUTABLE[config.device][1]
    cmdline += ["-a", f"{project_name(config)}/{output_dir(config)}/{project_name(config)}.{config.compiler.image_ext}"]
    return cmdline


@matrix_command()
def eventlist(config):
    return ["eventlist", "-s", "EventRecorder.log"]


if __name__ == "__main__":
    main()
