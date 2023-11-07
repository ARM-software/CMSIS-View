# Revision History {#rev_hist}

CMSIS-View version is offically updated upon releases of the [CMSIS-View pack](https://www.keil.arm.com/packs/cmsis-view-arm/versions/).

The table below provides information about the changes delivered with specific versions of CMSIS-View.

<table class="cmtable" summary="Revision History">
    <tr>
      <th>Version</th>
      <th>Description</th>
    </tr>
    <tr>
      <td>1.0.0</td>
      <td>
       Initial release of CMSIS-View with EventRecorder and Fault components as replacement for Keil.ARM-Compiler pack.
        - Renamed component class to CMSIS-View.
        - Fixes/additions for IAR Compiler.
        - Documentation enhacements.
        - Optimized Record Lock/Unlock in Event Recorder (using C11 atomics except for Cortex-M0).
        - Corrected timestamp overflow handling in Event Recorder.
        - Added "CMSIS-View::Fault" component for recording system faults.
       </td>
    </tr>

</table>
